package example

import (
	"fmt"
	"gonum.org/v1/gonum/stat/combin"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"testing"
)

func binomialDistribution(n, k int, p float64) float64 {
	// f(n,k,p) = (n choose k) p^k (1-p)^(n-k)
	return float64(combin.Binomial(n, k)) * math.Pow(p, float64(k)) * math.Pow(1.0-p, float64(n-k))
}

func TestPlotBinomialDistribution(t *testing.T) {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Probability of n successes in n trails, where probability of success p is"
	p.X.Label.Text = "n - number of trials"
	p.Y.Label.Text = "probability of success in n - trails"
	p.Legend.Top = true

	lines := make([]interface{}, 0)
	for i := 0; i <= 20; i++ {
		p := 0.05 * float64(i)
		var points plotter.XYs
		for n := 1; n <= 30; n++ {
			points = append(points, plotter.XY{
				X: float64(n),
				Y: binomialDistribution(n, n, p),
			})
		}

		label := fmt.Sprintf("p = %d%%", int(p*100))
		lines = append(lines, label, points)
	}

	err = plotutil.AddLinePoints(p, lines...)

	// Save the plot to a PNG file.
	if err := p.Save(18*vg.Inch, 9*vg.Inch, "binomial_distribution_test.png"); err != nil {
		panic(err)
	}
}
