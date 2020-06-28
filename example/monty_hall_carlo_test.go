package example

import (
	"math/rand"
	"testing"
)

// Probability of wining in Monty Hall game
func TestMontyHallCarlo_switch_gates(t *testing.T) {
	n := 100000
	nonSwitchSuccesses := 0
	switchSuccesses := 0
	gates := []int{1, 0, 0}
	decisions := []int{0, 1, 2}

	for i := 0; i < n; i++ {
		// New game, different gates
		rand.Shuffle(len(gates), func(i, j int) {
			gates[i], gates[j] = gates[j], gates[i]
		})
		// New player, new decisions
		rand.Shuffle(len(decisions), func(i, j int) {
			decisions[i], decisions[j] = decisions[j], decisions[i]
		})

		// Player select the gate, and sticks to decision
		selectedGate := decisions[0]
		if gates[selectedGate] == 1 {
			nonSwitchSuccesses++
		}

		// In case player decides to switch gates,
		// Let's reduce number of player choices
		remainingDecisions := decisions[1:]

		// Gate that always will be relived is not winning gate (assumption)
		// So let's find it
		reviledGate := remainingDecisions[0]
		switchToGate := remainingDecisions[1]
		if gates[reviledGate] == 1 {
			reviledGate = remainingDecisions[1]
			switchToGate = remainingDecisions[0]
		}

		// Let's count player success when switching the doors result in winning the game
		if gates[switchToGate] == 1 {
			switchSuccesses++
		}
	}

	between(t, float64(nonSwitchSuccesses)/float64(n), 0.32, 0.35, 0.001)
	between(t, float64(switchSuccesses)/float64(n), 0.64, 0.67, 0.001)
}
