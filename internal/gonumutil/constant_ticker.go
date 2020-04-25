package gonumutil

import (
	"fmt"
	"gonum.org/v1/plot"
)

func NewConstantNumTicker(step uint) plot.TickerFunc {
	return NewConstantTicker("%1.0f", float64(step))
}

func NewConstantTicker(label string, step float64) plot.TickerFunc {
	return plot.TickerFunc(func(min, max float64) []plot.Tick {
		var res = make([]plot.Tick, 0)

		val := min
		for val <= max {
			res = append(res, plot.Tick{
				Value: val,
				Label: fmt.Sprintf(label, val),
			})
			val += step
		}
		return res
	})
}
