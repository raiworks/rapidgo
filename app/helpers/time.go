package helpers

import (
	"fmt"
	"time"
)

// TimeAgo returns a human-readable relative time: "5 minutes ago".
func TimeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	default:
		return t.Format("Jan 2, 2006")
	}
}

// FormatDate formats time as "Jan 2, 2006 3:04 PM".
func FormatDate(t time.Time) string {
	return t.Format("Jan 2, 2006 3:04 PM")
}
