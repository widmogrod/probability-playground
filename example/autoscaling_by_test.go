package example

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

// percent is a value between 0 and 100
type percent = float64

// num is value grater than or equal 0
type num = float64

type Range struct {
	Min percent
	Max percent
}

func (ir Range) Contains(v percent) bool {
	if ir.Min > v || ir.Max < v {
		return false
	}

	return true
}

type Context struct {
	// CPUNoopRange defines boundaries in which CPU should not trigger auto-scaling
	CPUNoopRange Range
	// MaintainsCPUAvg level of CPU utilisation that should be maintained when decision about scaling up or down in made
	MaintainsCPUAvg percent
	// CPUUtilisation represents current CPU utilisation.
	CPUUtilisation percent
	// Instances represents current number of instances of a service.
	Instances num
}

// CPUScale calculates how many instances should added or removed to maintain given CPU utilization
func CPUScale(in Context) int {
	if in.CPUNoopRange.Contains(in.CPUUtilisation) {
		return 0
	}

	return ScaleInstances(in.Instances, in.CPUUtilisation, in.MaintainsCPUAvg)
}

// ScaleInstances calculates how many instances should be added or removed
// to maintain given percentage of utilization of resource with respect to current utilization.
// Utilisation is abstract, and it can be applied to average CPU utilisation, average queue size,...
func ScaleInstances(instances, utilisation, maintain percent) int {
	candidate := instances * utilisation / maintain
	candidate = math.Ceil(candidate - instances)
	return int(candidate)
}

type Recommendation struct {
	ScaleUp   uint
	ScaleDown uint
}

func toRecommendation(candidate int) Recommendation {
	result := Recommendation{}
	if candidate < 0 {
		result.ScaleDown = uint(-candidate)
	} else if candidate > 0 {
		result.ScaleUp = uint(candidate)
	}

	return result
}

func TestAutoScalingBy(t *testing.T) {
	useCases := map[string]struct {
		ctx      Context
		expected Recommendation
	}{
		"maintain": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  85,
				Instances:       3,
			},
			expected: Recommendation{
				ScaleUp:   0,
				ScaleDown: 0,
			},
		},
		"CPUScale up - small": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  91,
				Instances:       3,
			},
			expected: Recommendation{
				ScaleUp:   1,
				ScaleDown: 0,
			},
		},
		"CPUScale up - big": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  99,
				Instances:       30,
			},
			expected: Recommendation{
				ScaleUp:   5,
				ScaleDown: 0,
			},
		},
		"CPUScale down - small": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  67,
				Instances:       4,
			},
			expected: Recommendation{
				ScaleUp:   0,
				ScaleDown: 0,
			},
		},
		"CPUScale down - big": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  71,
				Instances:       33,
			},
			expected: Recommendation{
				ScaleUp:   0,
				ScaleDown: 5,
			},
		},
	}
	for name, uc := range useCases {
		t.Run(name, func(t *testing.T) {
			result := toRecommendation(CPUScale(uc.ctx))
			assert.Equal(t, uc.expected, result)
		})
	}
}
