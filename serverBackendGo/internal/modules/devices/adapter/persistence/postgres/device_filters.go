package postgres

import (
	"fmt"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
)

const (
	msTwoHours  = 2 * 3600 * 1000
	msFourHours = 4 * 3600 * 1000
)

const searchExtraJoins = `
	LEFT JOIN configurations cfg ON d.configurationid = cfg.id
`

func sortDirection(req domain.SearchRequest) string {
	if req.SortDir != nil && strings.EqualFold(strings.TrimSpace(*req.SortDir), "desc") {
		return "DESC"
	}
	return "ASC"
}

// orderExpr returns ORDER BY for the outer device row fetch (one row per id).
func orderExpr(req domain.SearchRequest) string {
	return orderExprInner(req, false)
}

// orderExprGrouped returns ORDER BY for GROUP BY d.id page subquery.
func orderExprGrouped(req domain.SearchRequest) string {
	return orderExprInner(req, true)
}

func orderExprInner(req domain.SearchRequest, grouped bool) string {
	dir := sortDirection(req)
	sortBy := ""
	if req.SortBy != nil {
		sortBy = strings.ToUpper(strings.TrimSpace(*req.SortBy))
	}
	agg := func(col string) string {
		if grouped {
			if dir == "DESC" {
				return "MAX(" + col + ")"
			}
			return "MIN(" + col + ")"
		}
		return col
	}
	delta := "EXTRACT(EPOCH FROM NOW()) * 1000 - d.lastupdate"
	switch sortBy {
	case "LAST_UPDATE":
		return fmt.Sprintf("%s %s NULLS LAST, lower(d.number) ASC", agg("d.lastupdate"), dir)
	case "NUMBER":
		return fmt.Sprintf("lower(d.number) %s, d.id ASC", dir)
	case "IMEI":
		return fmt.Sprintf("lower(COALESCE(d.imei, d.infojson->>'imei', '')) %s, d.id ASC", dir)
	case "PHONE":
		return fmt.Sprintf("lower(COALESCE(d.phone, d.infojson->>'phone', '')) %s, d.id ASC", dir)
	case "MODEL":
		return fmt.Sprintf("lower(COALESCE(d.infojson->>'model', '')) %s, d.id ASC", dir)
	case "ANDROID_VERSION":
		return fmt.Sprintf("COALESCE(d.infojson->>'androidVersion', '') %s, d.id ASC", dir)
	case "MDM_MODE":
		return fmt.Sprintf("COALESCE(d.infojson->>'mdmMode', '') %s, d.id ASC", dir)
	case "KIOSK_MODE":
		return fmt.Sprintf("COALESCE(d.infojson->>'kioskMode', '') %s, d.id ASC", dir)
	case "BATTERY_LEVEL":
		return fmt.Sprintf("LPAD(COALESCE(d.infojson->>'batteryLevel', '0'), 3, '0') %s, d.id ASC", dir)
	case "DESCRIPTION":
		return fmt.Sprintf("lower(COALESCE(d.description, '')) %s, d.id ASC", dir)
	case "CONFIGURATION":
		return fmt.Sprintf("lower(COALESCE(cfg.name, '')) %s, d.id ASC", dir)
	case "STATUS":
		return fmt.Sprintf(`CASE
			WHEN (%s) < %d THEN 1
			WHEN (%s) < %d THEN 2
			ELSE 3
		END %s, lower(d.number) ASC`, delta, msTwoHours, delta, msFourHours, dir)
	default:
		return "lower(d.number) ASC, d.id ASC"
	}
}

func searchFilters(req domain.SearchRequest, args *[]any, where *string, argN *int) {
	if req.Value != nil && strings.TrimSpace(*req.Value) != "" {
		fast := req.FastSearch != nil && *req.FastSearch
		if fast {
			*where += fmt.Sprintf(` AND (d.number = $%d OR d.fastsearch = $%d)`, *argN, *argN)
			*args = append(*args, strings.TrimSpace(*req.Value))
			*argN++
		} else {
			pat := *req.Value
			*where += fmt.Sprintf(` AND (
				d.number ILIKE $%d OR d.description ILIKE $%d OR d.imei ILIKE $%d OR d.phone ILIKE $%d
				OR d.publicip ILIKE $%d
				OR d.infojson->>'imei' ILIKE $%d OR d.infojson->>'phone' ILIKE $%d
				OR d.infojson->>'model' ILIKE $%d OR d.infojson->>'serial' ILIKE $%d
				OR d.custom1 ILIKE $%d OR d.custom2 ILIKE $%d OR d.custom3 ILIKE $%d
				OR d.oldnumber ILIKE $%d OR cfg.name ILIKE $%d OR g.name ILIKE $%d
			)`, *argN, *argN, *argN, *argN, *argN, *argN, *argN, *argN, *argN, *argN, *argN, *argN, *argN, *argN, *argN)
			for i := 0; i < 15; i++ {
				*args = append(*args, pat)
			}
			*argN += 15
		}
	}
	if req.GroupID != nil && *req.GroupID > 0 {
		*where += fmt.Sprintf(` AND g.id = $%d`, *argN)
		*args = append(*args, *req.GroupID)
		*argN++
	}
	if req.ConfigurationID != nil && *req.ConfigurationID > 0 {
		*where += fmt.Sprintf(` AND d.configurationid = $%d`, *argN)
		*args = append(*args, *req.ConfigurationID)
		*argN++
	}
	if req.Status != nil {
		appendStatusFilter(*req.Status, args, where, argN)
	}
	if req.DateFrom != nil && *req.DateFrom > 0 {
		*where += fmt.Sprintf(` AND d.lastupdate >= $%d`, *argN)
		*args = append(*args, *req.DateFrom)
		*argN++
	}
	if req.DateTo != nil && *req.DateTo > 0 {
		*where += fmt.Sprintf(` AND d.lastupdate <= $%d`, *argN)
		*args = append(*args, *req.DateTo)
		*argN++
	}
	if req.OnlineEarlierMillis != nil && *req.OnlineEarlierMillis > 0 {
		*where += fmt.Sprintf(` AND d.lastupdate <= EXTRACT(EPOCH FROM NOW()) * 1000 - $%d`, *argN)
		*args = append(*args, *req.OnlineEarlierMillis)
		*argN++
	}
	if req.OnlineLaterMillis != nil && *req.OnlineLaterMillis > 0 {
		*where += fmt.Sprintf(` AND d.lastupdate >= EXTRACT(EPOCH FROM NOW()) * 1000 - $%d`, *argN)
		*args = append(*args, *req.OnlineLaterMillis)
		*argN++
	}
	if req.EnrollmentDateFrom != nil && *req.EnrollmentDateFrom > 0 {
		*where += fmt.Sprintf(` AND d.enrolltime >= $%d`, *argN)
		*args = append(*args, *req.EnrollmentDateFrom)
		*argN++
	}
	if req.EnrollmentDateTo != nil && *req.EnrollmentDateTo > 0 {
		*where += fmt.Sprintf(` AND d.enrolltime <= $%d`, *argN)
		*args = append(*args, *req.EnrollmentDateTo)
		*argN++
	}
	if req.MdmMode != nil {
		*where += fmt.Sprintf(` AND (d.infojson->>'mdmMode')::BOOLEAN = $%d`, *argN)
		*args = append(*args, *req.MdmMode)
		*argN++
	}
	if req.KioskMode != nil {
		*where += fmt.Sprintf(` AND (d.infojson->>'kioskMode')::BOOLEAN = $%d`, *argN)
		*args = append(*args, *req.KioskMode)
		*argN++
	}
	if req.AndroidVersion != nil && strings.TrimSpace(*req.AndroidVersion) != "" {
		*where += fmt.Sprintf(` AND d.infojson->>'androidVersion' = $%d`, *argN)
		*args = append(*args, strings.TrimSpace(*req.AndroidVersion))
		*argN++
	}
	if req.LauncherVersion != nil && strings.TrimSpace(*req.LauncherVersion) != "" {
		*where += fmt.Sprintf(` AND COALESCE(d.infojson->>'launcherVersion', '') = $%d`, *argN)
		*args = append(*args, strings.TrimSpace(*req.LauncherVersion))
		*argN++
	}
	if req.InstallationStatus != nil && strings.TrimSpace(*req.InstallationStatus) != "" {
		*where += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM jsonb_array_elements(COALESCE(d.infojson->'applications', '[]'::jsonb)) elem
			WHERE elem->>'status' = $%d
		)`, *argN)
		*args = append(*args, strings.TrimSpace(*req.InstallationStatus))
		*argN++
	}
	if req.ImeiChanged != nil && *req.ImeiChanged {
		*where += ` AND d.imeiupdatets / 1000 > EXTRACT(EPOCH FROM (NOW() - INTERVAL '1 hour'))`
	}
}

func appendStatusFilter(status string, args *[]any, where *string, argN *int) {
	s := strings.ToLower(strings.TrimSpace(status))
	if s == "" {
		return
	}
	delta := "EXTRACT(EPOCH FROM NOW()) * 1000 - d.lastupdate"
	switch s {
	case "green":
		*where += fmt.Sprintf(` AND (%s) < $%d`, delta, *argN)
		*args = append(*args, msTwoHours)
		*argN++
	case "yellow":
		*where += fmt.Sprintf(` AND (%s) >= $%d AND (%s) < $%d`, delta, *argN, delta, *argN+1)
		*args = append(*args, msTwoHours, msFourHours)
		*argN += 2
	case "red":
		*where += fmt.Sprintf(` AND (%s) >= $%d`, delta, *argN)
		*args = append(*args, msFourHours)
		*argN++
	}
}
