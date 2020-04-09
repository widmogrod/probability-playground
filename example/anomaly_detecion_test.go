package example

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"testing"
)

func TestAnomalyDetection(t *testing.T) {
	// P(Anomaly | Evidence) = P(E|H) * P(H) / P(E|H) * P(H) + P(E|~H) * P(~H)
	//
	//          +--------+------------------+
	//          |        |                  |
	//          |        |                  |
	//          |        |                  |
	//          +--------|                  |
	//          |        |                  |
	//          |        |                  |
	//  P(E|H)  |        +------------------+
	//          |        |                  | P(E|~H)
	//          +--------+------------------+
	//             P(H)         P(~H)
	// From a sampled request, what is probability of anomaly?
	//
	//
	//
	//                                    [A1]
	//                                     ▇▇
	//                                     ▇▇
	//               ▇▇▇▇▇                 ▇▇
	//             ▇▇▇▇▇▇▇▇	               ▇▇                                [A3]
	//           ▇▇▇▇▇▇▇▇▇▇▇▇▇          ▇▇▇▇▇▇▇▇▇▇▇▇            [A2]       ▇  ▇  ▇
	//          ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇    [A4]
	//    ----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|---> [t]
	//
	//
	//    1. Detect if in bucket of time, there is anomaly...
	//       Knowing probability of anomaly in bucket of time,
	//    2. Score requests that have highest level of causing anomaly
	//
	//  What types of anomaly there are?
	//  - spikes [A1]
	//  - no regularity, when it should be [A2]
	//  - jitter [A3]
	//  - no data [A4]
	//
	//  What information we have in request?
	//  - IP
	//  - User-Agent
	//  - Method
	//  - Body
	//  - Referer
	//
	// TODO:
	// - create a vector of features [is_mobile, is_web, ....]
	// - inference probability of vector P(anomaly|<vector>)
	// - naive approach?
	// - how to train it/simulate concept by semi-supervised learning?
	//
	// Example vector can have
	// - authorisation page: 1
	// - user_me_page: 1
	// - question_page: 10
	// - is_android:
	// - is_mobile:
	// - is_web:
	// - is_mobile_web:
	//
	// Cluster vectors and take a look what you can find...

	rand.Seed(int64(0))
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Let's take a look at rate of change
	var points, change, change2, change3, average plotter.XYs
	for i := 0; i < 150; i++ {
		r := float64(i) * 0.1

		if i > 50 && i < 60 {
			r += rand.Float64()
		}

		s := math.Abs(math.Sin(r))
		//s *= math.Pow(s, 5)
		if i > 80 && i < 90 {
			s *= rand.Float64() * 3
		}

		points = append(points, plotter.XY{
			X: float64(i),
			Y: s,
		})

		delta := onWindow(2, points, func(ys ...float64) float64 {
			return ys[0] - ys[1]
		})

		prev2 := onWindow(3, points, func(ys ...float64) float64 {
			return ys[0] - ys[1]
		})

		change = append(change, plotter.XY{
			X: float64(i),
			Y: changeH(points),
		})

		change2 = append(change2, plotter.XY{
			X: float64(i),
			Y: delta - prev2,
		})

		change3 = append(change3, plotter.XY{
			X: float64(i),
			Y: math.Sqrt(math.Pow(delta-prev2, 2)),
		})

		if i >= 10 {
			average = append(average, plotter.XY{
				X: float64(i),
				Y: normalise(change3),
			})
		}
	}

	err = plotutil.AddLinePoints(p,
		"behaviour", points,
		//"change(1)", change,
		//"change(2)", change2,
		//"√(𝚫₁^2+𝚫₂^2)", change3,
		"average", average,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := p.Save(18*vg.Inch, 9*vg.Inch, "spam_filtering_bayes_test.png"); err != nil {
		t.Fatal(err)
	}
}

func avgH(windowSize int, xys plotter.XYs) float64 {
	i := len(xys) - 1
	avg := .0
	if i >= windowSize {
		agg := .0
		for j := 0; j < windowSize; j++ {
			agg += xys[i-j].Y
		}
		avg = agg / float64(windowSize)
	}

	return avg
}

func sumH(windowSize int, xys plotter.XYs) float64 {
	i := len(xys) - 1
	sum := .0
	if i >= windowSize {
		for j := 0; j < windowSize; j++ {
			sum += xys[i-j].Y
		}
	}

	return sum
}

func changeH(xys plotter.XYs) float64 {
	i := len(xys) - 1
	diff := .0
	if i > 1 {
		diff = xys[i].Y - xys[i-1].Y
	}

	return diff
}

func onWindow(windowSize int, xys plotter.XYs, fn func(ys ...float64) float64) float64 {
	i := len(xys)
	if i < windowSize {
		return 0
	}

	attr := make([]float64, windowSize+1)
	for j := windowSize; j > 0; j-- {
		attr[windowSize-j] = xys[i-j].Y
	}
	return fn(attr...)
}

func normalise(xys plotter.XYs) float64 {
	return onWindow(1, xys, func(ys ...float64) float64 {
		return 1 - (1 / (1 + ys[0]))
	})
}

func expAvg(windowSize int, xys plotter.XYs) float64 {
	return onWindow(windowSize, xys, func(ys ...float64) float64 {
		var res float64
		for _, y := range ys {
			res += (1 - (1 / (1 + y)))
		}

		return res / float64(windowSize)
	})
}
