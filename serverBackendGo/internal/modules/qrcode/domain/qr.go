package domain

// QRQuery holds enrollment query parameters.
type QRQuery struct {
	DeviceID       string
	CreateOnDemand string
	UseID          string
	Groups         []string
	Size           int
}
