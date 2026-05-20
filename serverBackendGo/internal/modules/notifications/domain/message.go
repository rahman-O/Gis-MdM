package domain

// PlainPushMessage is the agent-facing push payload.
type PlainPushMessage struct {
	ID          int64  `json:"id"`
	MessageType string `json:"messageType"`
	Payload     string `json:"payload"`
}
