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

func TestOrderExpr_lastUpdateDesc(t *testing.T) {
	sortBy := "LAST_UPDATE"
	sortDir := "desc"
	req := domain.SearchRequest{SortBy: &sortBy, SortDir: &sortDir}
	got := orderExpr(req)
	if !strings.Contains(got, "d.lastupdate DESC") {
		t.Fatalf("got %q", got)
	}
}
