package domain

// LocationRecord represents a single GPS data point stored in device_locations.
type LocationRecord struct {
	ID           int64    `json:"id"`
	DeviceID     int      `json:"deviceId"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Accuracy     float64  `json:"accuracy"`
	Speed        float64  `json:"speed"`
	Altitude     *float64 `json:"altitude,omitempty"`
	BatteryLevel *float64 `json:"batteryLevel,omitempty"`
	NetworkType  string   `json:"networkType"`
	TrackingMode string   `json:"trackingMode"`
	Timestamp    int64    `json:"timestamp"`
	ReceivedAt   int64    `json:"receivedAt,omitempty"`
	Month        string   `json:"month,omitempty"`
}

// LocationArchive represents an hourly summary of location data.
type LocationArchive struct {
	ID               int64   `json:"id"`
	DeviceID         int     `json:"deviceId"`
	HourStart        int64   `json:"hourStart"`
	StartLatitude    float64 `json:"startLatitude"`
	StartLongitude   float64 `json:"startLongitude"`
	EndLatitude      float64 `json:"endLatitude"`
	EndLongitude     float64 `json:"endLongitude"`
	DistanceTraveled float64 `json:"distanceTraveled"`
	PointCount       int     `json:"pointCount"`
}

// BatchUploadEntry is a single location record in a batch submission from the agent.
type BatchUploadEntry struct {
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Accuracy     float64  `json:"accuracy"`
	Speed        float64  `json:"speed"`
	Altitude     *float64 `json:"altitude,omitempty"`
	BatteryLevel *float64 `json:"batteryLevel,omitempty"`
	NetworkType  string   `json:"networkType"`
	TrackingMode string   `json:"trackingMode"`
	Timestamp    int64    `json:"timestamp"`
}

// BatchUploadResponse is returned after processing a batch of location records.
type BatchUploadResponse struct {
	Status   string   `json:"status"`
	Accepted int      `json:"accepted"`
	Rejected int      `json:"rejected"`
	Reasons  []string `json:"reasons,omitempty"`
}

// LocationUpdateData is broadcast via WebSocket when a new location is stored.
type LocationUpdateData struct {
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Accuracy     float64  `json:"accuracy"`
	Speed        float64  `json:"speed"`
	Altitude     *float64 `json:"altitude,omitempty"`
	BatteryLevel *float64 `json:"batteryLevel,omitempty"`
	NetworkType  string   `json:"networkType,omitempty"`
	Timestamp    int64    `json:"timestamp"`
}
