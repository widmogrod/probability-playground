package example

import (
	"math/rand"
	"testing"
)

// Amazon Mechanical Turk, provides capability to distribute simple task to humans to solve.
// This capability is also known under the name of crowd souring.
// Because is hard to verify quality of work in setup, the question arise:
// How probable is that collected work product (data) are like random answers?
//
// This test case aims to calculate probability of random answers where:
// - there is n-tasks to be solved by workers
// - each worker can make one of four decision per task (independent, with identical probability)
// - each tasks has to be answer k-times, by different workers
//
// To measure probability of agreement on decision measure of consensus is introduced
// - consensus of degree 3 - happens when three independent workers make the same decision on one task
// - consensus of degree 2 - happens when two independent workers choose the same decision, but one makes different
// - consensus of degree 1 - happens when each worker makes different decision
func TestAWSMechanicalTurkProbabilityOfConsensusMonteCarlo(t *testing.T) {
	n := 10000
	k := 3
	epsilon := 0.01
	consensusDegree3 := .0
	consensusDegree2 := .0
	consensusDegree1 := .0

	// there is n tasks to be solved by workers,
	for i := 0; i < n; i++ {
		voteA, voteB, voteC, voteD := 0, 0, 0, 0
		// each task must be answer k-times
		for w := 0; w < k; w++ {
			decision := rand.Float32()
			if decision <= 0.25 {
				voteA++
			} else if decision <= 0.5 {
				voteB++
			} else if decision <= 0.75 {
				voteC++
			} else {
				voteD++
			}
		}

		// when workers agree 3 times on the same decision we have consensus of degree 3
		if voteA >= 3 || voteB >= 3 || voteC >= 3 || voteD >= 3 {
			consensusDegree3++
		} else if voteA >= 2 || voteB >= 2 || voteC >= 2 || voteD >= 2 {
			// when workers agree 2 times on the same decision we have consensus of degree 2
			consensusDegree2++
		} else {
			// when each worker makes different decisions, then it's consensus of degree 1
			consensusDegree1++
		}
	}

	// What is probability of degree 3,2 and 1 consensus?
	// or in other words:
	// What is probability that in n-tasks workers making decision on random, will reach consensus on the same decision?
	degree3Probability := consensusDegree3 / float64(n)
	degree2Probability := consensusDegree2 / float64(n)
	degree1Probability := consensusDegree1 / float64(n)

	between(t, degree3Probability, 0.06, 0.08, epsilon)
	between(t, degree2Probability, 0.53, 0.55, epsilon)
	between(t, degree1Probability, 0.35, 0.38, epsilon)
}
