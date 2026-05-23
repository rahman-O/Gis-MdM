package application

import (
	"strconv"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/domain"
)

// MatchesSchedule returns true if task cron masks match time t (Java PushScheduleDAO.findMatchingTime).
func MatchesSchedule(task domain.PluginPushSchedule, t time.Time) bool {
	minBit := parseScheduleMask(task.Min, 60, false)
	hourBit := parseScheduleMask(task.Hour, 24, false)
	dayBit := parseScheduleMask(task.Day, 31, true)
	weekBit := parseScheduleMask(task.Weekday, 7, true)
	monthBit := parseScheduleMask(task.Month, 12, false)
	if minBit == "" || hourBit == "" || dayBit == "" || weekBit == "" || monthBit == "" {
		return false
	}
	minute := t.Minute()
	hour := t.Hour()
	day := t.Day()
	weekday := int(t.Weekday()) // Sunday=0 in Go; Java Calendar.DAY_OF_WEEK Sunday=1
	if weekday == 0 {
		weekday = 7
	}
	month := int(t.Month()) - 1
	return bitAt(minBit, minute) && bitAt(hourBit, hour) && bitAt(dayBit, day-1) &&
		bitAt(weekBit, weekday-1) && bitAt(monthBit, month)
}

func bitAt(mask string, idx int) bool {
	if idx < 0 || idx >= len(mask) {
		return false
	}
	return mask[idx] == '1'
}

func parseScheduleMask(raw string, length int, startFromOne bool) string {
	mask := strings.ReplaceAll(strings.TrimSpace(raw), " ", "")
	if mask == "" {
		mask = "*"
	}
	res := make([]byte, length)
	if mask == "*" {
		for i := range res {
			res[i] = '1'
		}
		return string(res)
	}
	if strings.HasPrefix(mask, "*/") {
		n, err := strconv.Atoi(mask[2:])
		if err != nil || n < 1 {
			return ""
		}
		for i := 0; i < length; i++ {
			if i%n == 0 {
				res[i] = '1'
			} else {
				res[i] = '0'
			}
		}
		return string(res)
	}
	n, err := strconv.Atoi(mask)
	if err != nil {
		return ""
	}
	if startFromOne {
		n--
	}
	if n < 0 || n >= length {
		return ""
	}
	for i := range res {
		if i == n {
			res[i] = '1'
		} else {
			res[i] = '0'
		}
	}
	return string(res)
}
