package example

import (
	"math"
	"testing"
)

func between(t *testing.T, value, min, max, epsilon float64) {
	lower := math.Max(value, value+epsilon)
	upper := math.Min(value, value-epsilon)
	if lower >= min && upper <= max {
		t.Logf("value between (+-%f):\n\t %f < %f <  %f", epsilon, min, value, max)
	} else {
		t.Errorf("value not between(+-%f):\n\t %f < %f <  %f", epsilon, min, value, max)
	}
}
