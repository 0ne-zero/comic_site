package route

import (
	"fmt"
	"strings"
	"time"
)

func HowManyAgo(t time.Time) string {
	hours := int(time.Since(t).Hours())
	if hours <= 23 {
		return fmt.Sprintf("%d ساعت پیش", hours)
	} else if hours <= 8760 {
		return fmt.Sprintf("%d روز پیش", int(hours/24))
	} else {
		return fmt.Sprintf("%d سال پیش", int((hours/24)/365))
	}
}
func IsEpisodeNew(t time.Time) bool {
	if int(time.Since(t).Hours()) > 72 {
		return false
	} else {
		return true
	}
}

func Minus(a, b int) int {
	return a - b
}

func Plus(a, b int) int {
	return a + b
}
func IsGetParameterExists(url string) bool {
	if _, _, found := strings.Cut(url, "?"); found {
		return true
	} else {
		return false
	}
}
