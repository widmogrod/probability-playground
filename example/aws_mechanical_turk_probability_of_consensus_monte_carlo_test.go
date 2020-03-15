package example

import (
	"golang.org/x/exp/errors/fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math/rand"
	"strconv"
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
// To measure the probability of agreement on decision let's introduce the concept called the degree of consensus:
// - consensus of degree 3 - happens when three independent workers make the same decision on one task
// - consensus of degree 2 - happens when two independent workers choose the same decision, but one makes different
// - consensus of degree 1 - happens when each worker makes different decision
func TestAWSMechanicalTurkProbabilityOfConsensusMonteCarlo(t *testing.T) {
	n := 10000
	k := 3
	epsilon := 0.01

	consensusProbability := mTurkMonteCarlo(n, k)

	between(t, consensusProbability[3], 0.06, 0.08, epsilon)
	between(t, consensusProbability[2], 0.53, 0.55, epsilon)
	between(t, consensusProbability[1], 0.35, 0.38, epsilon)
}

// mTurkMonteCarlo simulation where
// - n-tasks to be solved by workers
// - each tasks has to be answer k-times
func mTurkMonteCarlo(n, k int) []float64 {
	consensusDegrees := make([]float64, k+1)
	consensusProbability := make([]float64, k+1)

	// there is n tasks to be solved by workers,
	for i := 0; i < n; i++ {
		votes := make([]float64, 4)
		// each task must be answer k-times
		for w := 0; w < k; w++ {
			decision := rand.Float32()
			if decision <= 0.25 {
				votes[0]++
			} else if decision <= 0.5 {
				votes[1]++
			} else if decision <= 0.75 {
				votes[2]++
			} else {
				votes[3]++
			}
		}

		consensusReached := int(max(votes))
		consensusDegrees[consensusReached]++
	}

	// What is probability of degree 3,2 and 1 consensus?
	// or in other words:
	// What is probability that in n-tasks workers making decision on random, will reach consensus on the same decision?
	for degree, reachTimes := range consensusDegrees {
		consensusProbability[degree] = reachTimes / float64(n)
	}

	return consensusProbability
}

func max(xs []float64) float64 {
	result := .0
	for _, x := range xs {
		if x > result {
			result = x
		}
	}

	return result
}

func TestPlotDistributionOfAWSMechanicalTurkProbabilityOfConsensusMonteCarlo(t *testing.T) {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "mTurn degrees of consensus"
	p.X.Label.Text = "decisions per task"
	p.Y.Label.Text = "probability of degree of consensus"
	p.Legend.Top = true

	n := 20000
	maxWorkers := 37

	degrees := make([]interface{}, maxWorkers*2)
	for workers := 0; workers < maxWorkers; workers++ {
		degrees[workers*2] = fmt.Sprintf("Degree %d", workers)
		degrees[workers*2+1] = make(plotter.XYs, maxWorkers)
	}

	for workers := 0; workers < maxWorkers; workers++ {
		consensusProbability := mTurkMonteCarlo(n, workers)
		for degree, probability := range consensusProbability {
			if d, ok := degrees[degree*2+1].(plotter.XYs); ok {
				degrees[degree*2+1] = append(
					d,
					plotter.XY{float64(workers), probability},
				)
			}
		}
	}

	err = plotutil.AddLinePoints(p, degrees...)
	if err != nil {
		panic(err)
	}

	labels := []string{}
	for i := 0; i < maxWorkers; i++ {
		labels = append(labels, strconv.Itoa(i))
	}

	p.NominalX(labels...)

	// Save the plot to a PNG file.
	if err := p.Save(18*vg.Inch, 9*vg.Inch, "aws_mechanical_turk_probability_of_consensus_monte_carlo_test.png"); err != nil {
		panic(err)
	}
}
