package application

import (
	"testing"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/domain"
)

func TestMatchesSchedule_wildcard(t *testing.T) {
	task := domain.PluginPushSchedule{
		Min: "*", Hour: "*", Day: "*", Weekday: "*", Month: "*",
	}
	if !MatchesSchedule(task, time.Date(2026, 5, 21, 12, 30, 0, 0, time.UTC)) {
		t.Fatal("expected match for all wildcards")
	}
}

func TestMatchesSchedule_specificMinute(t *testing.T) {
	task := domain.PluginPushSchedule{
		Min: "30", Hour: "*", Day: "*", Weekday: "*", Month: "*",
	}
	ok := MatchesSchedule(task, time.Date(2026, 5, 21, 12, 30, 0, 0, time.UTC))
	if !ok {
		t.Fatal("expected minute 30 match")
	}
	no := MatchesSchedule(task, time.Date(2026, 5, 21, 12, 31, 0, 0, time.UTC))
	if no {
		t.Fatal("expected no match at minute 31")
	}
}
