package domain

// PushRequest mirrors React pushService / PushApiResource.
type PushRequest struct {
	MessageType   string   `json:"messageType"`
	Payload       string   `json:"payload"`
	DeviceNumbers []string `json:"deviceNumbers"`
	Groups        []string `json:"groups"`
	Broadcast     bool     `json:"broadcast"`
}
