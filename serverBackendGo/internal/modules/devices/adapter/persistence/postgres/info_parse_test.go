package postgres

import (
	"database/sql"
	"testing"
)

func TestParseDeviceInfo_fromInfoColumn(t *testing.T) {
	info := sql.NullString{
		String: `{"batteryLevel":85,"model":"Pixel","androidVersion":"14","mdmMode":true}`,
		Valid:  true,
	}
	out := parseDeviceInfo(info, nil, sql.NullInt64{}, sql.NullString{})
	if out == nil {
		t.Fatal("expected info")
	}
	if out.BatteryLevel == nil || *out.BatteryLevel != 85 {
		t.Fatalf("battery %v", out.BatteryLevel)
	}
	if out.Model == nil || *out.Model != "Pixel" {
		t.Fatalf("model %v", out.Model)
	}
}

func TestParseDeviceInfo_mergesInfojson(t *testing.T) {
	infojson := []byte(`{"batteryLevel":90,"kioskMode":false}`)
	out := parseDeviceInfo(sql.NullString{}, infojson, sql.NullInt64{}, sql.NullString{})
	if out == nil || out.BatteryLevel == nil || *out.BatteryLevel != 90 {
		t.Fatalf("got %+v", out)
	}
}
