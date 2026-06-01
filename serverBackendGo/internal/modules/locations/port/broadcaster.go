package port

import "github.com/gis-mdm/server-backend-go/internal/modules/locations/domain"

// Broadcaster defines the interface for real-time location broadcasting via WebSocket.
type Broadcaster interface {
	BroadcastLocation(deviceID string, data domain.LocationUpdateData)
	BroadcastDeviceStatus(deviceID string, status string)
	UpdateDeviceHeartbeat(deviceID string)
}
