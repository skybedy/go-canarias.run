package utils

import (
	"fmt"
	"strings"
)

// FormatDate converts DD/MM/YYYY strings to YYYY-MM-DD
func FormatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	// Try DD/MM/YYYY
	parts := strings.Split(dateStr, "/")
	if len(parts) == 3 {
		day := parts[0]
		month := parts[1]
		year := parts[2]
		if len(day) == 1 {
			day = "0" + day
		}
		if len(month) == 1 {
			month = "0" + month
		}
		return fmt.Sprintf("%s-%s-%s", year, month, day)
	}

	// Try DD.MM.YYYY
	parts = strings.Split(dateStr, ".")
	if len(parts) == 3 {
		day := strings.TrimSpace(parts[0])
		month := strings.TrimSpace(parts[1])
		year := strings.TrimSpace(parts[2])
		if len(day) == 1 {
			day = "0" + day
		}
		if len(month) == 1 {
			month = "0" + month
		}
		return fmt.Sprintf("%s-%s-%s", year, month, day)
	}

	return ""
}
