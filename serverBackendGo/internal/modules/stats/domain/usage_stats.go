package domain

// UsageStats mirrors com.hmdm.persistence.domain.UsageStats (PUT /rest/public/stats).
type UsageStats struct {
	InstanceID     string `json:"instanceId"`
	WebVersion     string `json:"webVersion,omitempty"`
	Community      bool   `json:"community"`
	DevicesTotal   int    `json:"devicesTotal"`
	DevicesOnline  int    `json:"devicesOnline"`
	CPUTotal       int    `json:"cpuTotal"`
	CPUUsed        int    `json:"cpuUsed"`
	RAMTotal       int    `json:"ramTotal"`
	RAMUsed        int    `json:"ramUsed"`
	Scheme         string `json:"scheme,omitempty"`
	Arch           string `json:"arch,omitempty"`
	OS             string `json:"os,omitempty"`
}
