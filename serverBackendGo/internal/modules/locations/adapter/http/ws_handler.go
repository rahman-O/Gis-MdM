package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/port"
)

const (
	wsWriteWait         = 10 * time.Second
	wsPingInterval      = 30 * time.Second
	wsPongWait          = 35 * time.Second
	wsMaxMessageSize    = 4096
	wsMaxSubscriptions  = 50
	deviceHeartbeatTTL  = 60 * time.Second
)

// ServerMessage is the JSON envelope sent from server to client.
type ServerMessage struct {
	Type     string      `json:"type"`
	DeviceID string      `json:"deviceId,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Status   string      `json:"status,omitempty"`
	Message  string      `json:"message,omitempty"`
}

// clientMessage is the JSON envelope received from client.
type clientMessage struct {
	Type     string `json:"type"`
	DeviceID string `json:"deviceId"`
}

// wsClient represents a single WebSocket connection.
type wsClient struct {
	conn          *websocket.Conn
	subscriptions map[string]bool
	mu            sync.Mutex
}

// WebSocketServer manages WebSocket connections and broadcasts location updates.
type WebSocketServer struct {
	upgrader websocket.Upgrader
	log      *slog.Logger

	// deviceID → set of clients subscribed to that device
	subscribers   map[string]map[*wsClient]bool
	subscribersMu sync.RWMutex

	// device heartbeat tracking
	heartbeats   map[string]time.Time
	heartbeatsMu sync.RWMutex

	heartbeatInterval time.Duration
	deviceTimeout     time.Duration
}

// NewWebSocketServer creates a new WebSocket server.
func NewWebSocketServer(log *slog.Logger, heartbeatSec, deviceTimeoutSec int) *WebSocketServer {
	if heartbeatSec <= 0 {
		heartbeatSec = 30
	}
	if deviceTimeoutSec <= 0 {
		deviceTimeoutSec = 60
	}
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development; restrict in production via reverse proxy
			},
		},
		log:               log,
		subscribers:       make(map[string]map[*wsClient]bool),
		heartbeats:        make(map[string]time.Time),
		heartbeatInterval: time.Duration(heartbeatSec) * time.Second,
		deviceTimeout:     time.Duration(deviceTimeoutSec) * time.Second,
	}
}

// Ensure WebSocketServer implements port.Broadcaster.
var _ port.Broadcaster = (*WebSocketServer)(nil)

// HandleUpgrade upgrades an HTTP connection to WebSocket and manages the client lifecycle.
func (ws *WebSocketServer) HandleUpgrade(c *gin.Context) {
	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		ws.log.Warn("websocket upgrade failed", "err", err, "remoteAddr", c.ClientIP())
		return
	}

	client := &wsClient{
		conn:          conn,
		subscriptions: make(map[string]bool),
	}

	// Auto-subscribe to device from URL path if present
	deviceID := c.Param("deviceId")
	if deviceID != "" {
		ws.subscribe(client, deviceID)
	}

	// Start read/write pumps
	go ws.readPump(client)
	go ws.writePing(client)
}

// readPump reads messages from the client and handles subscribe/unsubscribe.
func (ws *WebSocketServer) readPump(client *wsClient) {
	defer ws.removeClient(client)

	client.conn.SetReadLimit(wsMaxMessageSize)
	_ = client.conn.SetReadDeadline(time.Now().Add(wsPongWait))
	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				ws.log.Debug("websocket read error", "err", err)
			}
			return
		}

		var msg clientMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			ws.sendError(client, "invalid JSON message")
			continue
		}

		switch msg.Type {
		case "subscribe":
			if msg.DeviceID == "" {
				ws.sendError(client, "deviceId required for subscribe")
				continue
			}
			ws.subscribe(client, msg.DeviceID)
		case "unsubscribe":
			if msg.DeviceID == "" {
				ws.sendError(client, "deviceId required for unsubscribe")
				continue
			}
			ws.unsubscribe(client, msg.DeviceID)
		default:
			ws.sendError(client, "unknown message type: "+msg.Type)
		}
	}
}

// writePing sends periodic pings to keep the connection alive.
func (ws *WebSocketServer) writePing(client *wsClient) {
	ticker := time.NewTicker(wsPingInterval)
	defer ticker.Stop()

	for range ticker.C {
		client.mu.Lock()
		_ = client.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
		err := client.conn.WriteMessage(websocket.PingMessage, nil)
		client.mu.Unlock()
		if err != nil {
			return
		}
	}
}

// subscribe adds a client to a device's subscriber list.
func (ws *WebSocketServer) subscribe(client *wsClient, deviceID string) {
	client.mu.Lock()
	if len(client.subscriptions) >= wsMaxSubscriptions {
		client.mu.Unlock()
		ws.sendError(client, "max subscriptions reached (50)")
		return
	}
	if client.subscriptions[deviceID] {
		client.mu.Unlock()
		return // already subscribed
	}
	client.subscriptions[deviceID] = true
	client.mu.Unlock()

	ws.subscribersMu.Lock()
	if ws.subscribers[deviceID] == nil {
		ws.subscribers[deviceID] = make(map[*wsClient]bool)
	}
	ws.subscribers[deviceID][client] = true
	ws.subscribersMu.Unlock()

	ws.sendJSON(client, ServerMessage{Type: "subscribed", DeviceID: deviceID})
}

// unsubscribe removes a client from a device's subscriber list.
func (ws *WebSocketServer) unsubscribe(client *wsClient, deviceID string) {
	client.mu.Lock()
	delete(client.subscriptions, deviceID)
	client.mu.Unlock()

	ws.subscribersMu.Lock()
	if subs, ok := ws.subscribers[deviceID]; ok {
		delete(subs, client)
		if len(subs) == 0 {
			delete(ws.subscribers, deviceID)
		}
	}
	ws.subscribersMu.Unlock()

	ws.sendJSON(client, ServerMessage{Type: "unsubscribed", DeviceID: deviceID})
}

// removeClient cleans up all subscriptions and closes the connection.
func (ws *WebSocketServer) removeClient(client *wsClient) {
	client.mu.Lock()
	subs := make([]string, 0, len(client.subscriptions))
	for deviceID := range client.subscriptions {
		subs = append(subs, deviceID)
	}
	client.subscriptions = nil
	client.mu.Unlock()

	ws.subscribersMu.Lock()
	for _, deviceID := range subs {
		if clients, ok := ws.subscribers[deviceID]; ok {
			delete(clients, client)
			if len(clients) == 0 {
				delete(ws.subscribers, deviceID)
			}
		}
	}
	ws.subscribersMu.Unlock()

	_ = client.conn.Close()
}

// BroadcastLocation sends a location_update message to all subscribers of a device.
func (ws *WebSocketServer) BroadcastLocation(deviceID string, data domain.LocationUpdateData) {
	msg := ServerMessage{
		Type:     "location_update",
		DeviceID: deviceID,
		Data:     data,
	}

	ws.subscribersMu.RLock()
	clients := ws.subscribers[deviceID]
	// Copy to avoid holding lock during writes
	targets := make([]*wsClient, 0, len(clients))
	for client := range clients {
		targets = append(targets, client)
	}
	ws.subscribersMu.RUnlock()

	for _, client := range targets {
		ws.sendJSON(client, msg)
	}
}

// BroadcastDeviceStatus sends an online/offline status message to all subscribers.
func (ws *WebSocketServer) BroadcastDeviceStatus(deviceID string, status string) {
	msg := ServerMessage{
		Type:     "device_status",
		DeviceID: deviceID,
		Status:   status,
	}

	ws.subscribersMu.RLock()
	clients := ws.subscribers[deviceID]
	targets := make([]*wsClient, 0, len(clients))
	for client := range clients {
		targets = append(targets, client)
	}
	ws.subscribersMu.RUnlock()

	for _, client := range targets {
		ws.sendJSON(client, msg)
	}
}

// UpdateDeviceHeartbeat records the last activity time for a device.
func (ws *WebSocketServer) UpdateDeviceHeartbeat(deviceID string) {
	ws.heartbeatsMu.Lock()
	ws.heartbeats[deviceID] = time.Now()
	ws.heartbeatsMu.Unlock()

	// Broadcast online status
	ws.BroadcastDeviceStatus(deviceID, "online")
}

// StartHeartbeatMonitor checks for devices that have gone offline.
func (ws *WebSocketServer) StartHeartbeatMonitor(ctx context.Context) {
	ticker := time.NewTicker(ws.deviceTimeout / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ws.checkHeartbeats()
		}
	}
}

// checkHeartbeats detects devices that have exceeded the heartbeat timeout.
func (ws *WebSocketServer) checkHeartbeats() {
	now := time.Now()

	ws.heartbeatsMu.Lock()
	var offline []string
	for deviceID, lastSeen := range ws.heartbeats {
		if now.Sub(lastSeen) > ws.deviceTimeout {
			offline = append(offline, deviceID)
			delete(ws.heartbeats, deviceID)
		}
	}
	ws.heartbeatsMu.Unlock()

	for _, deviceID := range offline {
		ws.BroadcastDeviceStatus(deviceID, "offline")
	}
}

// Close terminates all active WebSocket connections.
func (ws *WebSocketServer) Close() {
	ws.subscribersMu.Lock()
	defer ws.subscribersMu.Unlock()

	for _, clients := range ws.subscribers {
		for client := range clients {
			_ = client.conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseGoingAway, "server shutting down"),
				time.Now().Add(wsWriteWait),
			)
			_ = client.conn.Close()
		}
	}
	ws.subscribers = make(map[string]map[*wsClient]bool)
}

// sendJSON writes a JSON message to a client.
func (ws *WebSocketServer) sendJSON(client *wsClient, msg ServerMessage) {
	client.mu.Lock()
	defer client.mu.Unlock()

	_ = client.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
	if err := client.conn.WriteJSON(msg); err != nil {
		ws.log.Debug("websocket write failed", "err", err)
	}
}

// sendError sends an error message to a client.
func (ws *WebSocketServer) sendError(client *wsClient, message string) {
	ws.sendJSON(client, ServerMessage{Type: "error", Message: message})
}
