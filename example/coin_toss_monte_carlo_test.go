package example

import (
	"math"
	"math/rand"
	"testing"
)

func TestCoinTossMonteCarlo(t *testing.T) {
	n := 1000
	epsilon := 0.05
	heads, tails := .0, .0
	for i := 0; i < n; i++ {
		if rand.Float32() <= 0.5 {
			heads++
		} else {
			tails++
		}
	}

	headProbability := heads / float64(n)
	tailProbability := tails / float64(n)

	between(t, headProbability, 0, 0.5, epsilon)
	between(t, tailProbability, 0, 0.5, epsilon)
}

func between(t *testing.T, value, min, max, epsilon float64) {
	lower := math.Max(value, value+epsilon)
	upper := math.Min(value, value-epsilon)
	if lower >= min && upper <= max {
		t.Logf("value between (+-%f):\n\t %f < %f <  %f", epsilon, min, value, max)
	} else {
		t.Errorf("value not between(+-%f):\n\t %f < %f <  %f", epsilon, min, value, max)
	}
}
