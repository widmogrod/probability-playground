package example

import (
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
