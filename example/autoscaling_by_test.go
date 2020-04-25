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

type InRange struct {
	Min percent
	Max percent
}

func (ir InRange) Contains(v percent) bool {
	if ir.Min > v || ir.Max < v {
		return false
	}

	return true
}

type Context struct {
	CPUInRange      InRange
	MaintainsCPUAvg percent
	CPUUtilisation  percent
	Instances       num
}

// CPUScale calculates how many instances should added or removed to maintain given CPU utilization
func CPUScale(in Context) float64 {
	if in.CPUUtilisation <= in.CPUInRange.Max && in.CPUUtilisation >= in.CPUInRange.Min {
		return 0
	}

	candidate := in.Instances * in.CPUUtilisation / in.MaintainsCPUAvg
	candidate = math.Ceil(candidate - in.Instances)
	return candidate
}

type Recommendation struct {
	ScaleUp   num
	ScaleDown num
}

func toRecommendation(candidate float64) Recommendation {
	result := Recommendation{}
	if candidate < 0 {
		result.ScaleDown = -candidate
	} else if candidate > 0 {
		result.ScaleUp = candidate
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
				CPUInRange: InRange{
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
				CPUInRange: InRange{
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
				CPUInRange: InRange{
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
				CPUInRange: InRange{
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
				CPUInRange: InRange{
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
