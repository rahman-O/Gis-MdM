package postgres

import (
	"strings"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
)

func TestSearchFilters_statusGreen(t *testing.T) {
	status := "green"
	req := domain.SearchRequest{Status: &status}
	var args []any
	where := "WHERE 1=1"
	argN := 1
	searchFilters(req, &args, &where, &argN)
	if !strings.Contains(where, "EXTRACT(EPOCH FROM NOW())") {
		t.Fatalf("expected lastupdate band filter, got %q", where)
	}
	if len(args) != 1 || args[0] != msTwoHours {
		t.Fatalf("args %v", args)
	}
}

func TestSearchFilters_statusYellow(t *testing.T) {
	status := "yellow"
	req := domain.SearchRequest{Status: &status}
	var args []any
	where := "WHERE 1=1"
	argN := 1
	searchFilters(req, &args, &where, &argN)
	if len(args) != 2 || args[0] != msTwoHours || args[1] != msFourHours {
		t.Fatalf("args %v", args)
	}
}

func TestSearchFilters_installationStatusUsesDeviceStatuses(t *testing.T) {
	status := "FAILURE"
	req := domain.SearchRequest{InstallationStatus: &status}
	var args []any
	where := "WHERE 1=1"
	argN := 1
	searchFilters(req, &args, &where, &argN)
	if !strings.Contains(where, "ds.applicationsstatus") {
		t.Fatalf("expected devicestatuses filter, got %q", where)
	}
	if len(args) != 1 || args[0] != "FAILURE" {
		t.Fatalf("args %v", args)
	}
}

func TestNeedsDeviceStatusJoin(t *testing.T) {
	sortBy := "INSTALLATIONS"
	if !needsDeviceStatusJoin(domain.SearchRequest{SortBy: &sortBy}) {
		t.Fatal("INSTALLATIONS sort should need devicestatuses join")
	}
	inst := "SUCCESS"
	if !needsDeviceStatusJoin(domain.SearchRequest{InstallationStatus: &inst}) {
		t.Fatal("installationStatus filter should need join")
	}
	if needsDeviceStatusJoin(domain.SearchRequest{}) {
		t.Fatal("empty request should not need join")
	}
}

func TestOrderExpr_installations(t *testing.T) {
	sortBy := "INSTALLATIONS"
	req := domain.SearchRequest{SortBy: &sortBy}
	got := orderExpr(req)
	if !strings.Contains(got, "ds.applicationsstatus") {
		t.Fatalf("got %q", got)
	}
}

func TestOrderExpr_lastUpdateDesc(t *testing.T) {
	sortBy := "LAST_UPDATE"
	sortDir := "desc"
	req := domain.SearchRequest{SortBy: &sortBy, SortDir: &sortDir}
	got := orderExpr(req)
	if !strings.Contains(got, "d.lastupdate DESC") {
		t.Fatalf("got %q", got)
	}
}
