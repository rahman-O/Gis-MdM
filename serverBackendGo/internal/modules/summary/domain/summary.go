package domain

import (
	"fmt"
	"time"
)

// ChartItem mirrors com.hmdm.rest.json.ChartItem.
type ChartItem struct {
	StringAttr string `json:"stringAttr,omitempty"`
	IntAttr    int    `json:"intAttr,omitempty"`
	Number     int    `json:"number,omitempty"`
}

// DeviceStats mirrors com.hmdm.rest.json.SummaryResponse (dashboard).
type DeviceStats struct {
	StatusSummary            []ChartItem `json:"statusSummary"`
	InstallSummary           []ChartItem `json:"installSummary"`
	DevicesTotal             int64       `json:"devicesTotal"`
	DevicesEnrolled          int64       `json:"devicesEnrolled"`
	DevicesEnrolledLastMonth int64       `json:"devicesEnrolledLastMonth"`
	DevicesEnrolledMonthly   []ChartItem `json:"devicesEnrolledMonthly"`
	TopConfigs               []string    `json:"topConfigs"`
	StatusOfflineByConfig    []int       `json:"statusOfflineByConfig"`
	StatusIdleByConfig       []int       `json:"statusIdleByConfig"`
	StatusOnlineByConfig     []int       `json:"statusOnlineByConfig"`
	AppFailureByConfig       []int       `json:"appFailureByConfig"`
	AppMismatchByConfig      []int       `json:"appMismatchByConfig"`
	AppSuccessByConfig       []int       `json:"appSuccessByConfig"`
}

// EmptyDeviceStats returns zeroed stats matching Java shape (used until devices table is migrated).
func EmptyDeviceStats() *DeviceStats {
	return &DeviceStats{
		StatusSummary: []ChartItem{
			{StringAttr: "green", Number: 0},
			{StringAttr: "yellow", Number: 0},
			{StringAttr: "red", Number: 0},
		},
		InstallSummary: []ChartItem{
			{StringAttr: "SUCCESS", Number: 0},
			{StringAttr: "VERSION_MISMATCH", Number: 0},
			{StringAttr: "FAILURE", Number: 0},
		},
		DevicesEnrolledMonthly: monthlyEnrollmentLabels(time.Now()),
		TopConfigs:             []string{},
		StatusOfflineByConfig:  []int{},
		StatusIdleByConfig:     []int{},
		StatusOnlineByConfig:   []int{},
		AppFailureByConfig:     []int{},
		AppMismatchByConfig:    []int{},
		AppSuccessByConfig:     []int{},
	}
}

func monthlyEnrollmentLabels(now time.Time) []ChartItem {
	year := now.Year() - 1
	month := int(now.Month())
	if month >= 12 {
		month = 0
		year++
	}
	items := make([]ChartItem, 0, 12)
	m, y := month, year
	for i := 0; i < 12; i++ {
		items = append(items, ChartItem{
			StringAttr: formatMonthLabel(m+1, y%100),
			Number:     0,
		})
		m++
		if m >= 12 {
			m = 0
			y++
		}
	}
	return items
}

func formatMonthLabel(month, year int) string {
	return fmt.Sprintf("%02d/%02d", month, year)
}
