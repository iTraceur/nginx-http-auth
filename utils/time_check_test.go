package utils

import (
	"testing"
)

func TestTimeCheck(t *testing.T) {
	timeRanges := []string{"11:02-17:11", "05:00"}
	b := TimeCheck(timeRanges)
	t.Log(b)
}
