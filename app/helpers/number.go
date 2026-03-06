package helpers

import (
	"fmt"
	"math"
)

// FormatBytes converts bytes to human-readable: 1536 → "1.50 KB".
func FormatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	if bytes == 0 {
		return "0 B"
	}
	i := int(math.Log(float64(bytes)) / math.Log(1024))
	if i >= len(units) {
		i = len(units) - 1
	}
	return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(1024, float64(i)), units[i])
}

// Clamp restricts a value between min and max.
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
