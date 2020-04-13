package example

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"testing"
)

// Simulation of Birthday Problem using Monte Carlo Method,
// shows how easy it is to answer question:
//
// > How probable is that k-people will share the same birthday?
//
// This question by itself is not fascinating,
// but intuition that is build after understand it is.
//
// You can read more in the Internet, what I want to highlight is that
// birthday paradox help to answer question what is probability of random hash collision.
// Important to have in back of your mind when building systems that are based on checksums, and randomly generated hashes.
// In short, such collision is more probable that you can think of ;)
func birthdayProblemMonteCarlo(samples, n, k int) float64 {
	success := .0
	// for each samples group of k participants
	for i := 0; i < samples; i++ {
		peopleWithBirthday := make([]bool, n)
		// ask each participant
		for w := 0; w < k; w++ {
			// for birth date
			day := rand.Intn(n)
			// and when at least two share the same date
			if peopleWithBirthday[day] {
				// count group as successful example
				success++
				break
			} else {
				peopleWithBirthday[day] = true
			}
		}
	}

	return success / float64(samples)
}

// Theoretical probability calculated following complement rule.
// Given group of k-people probability of at least two people sharing birthday
// in the same as probability of non of k-people sharing birthday, subtracted from one.
//
// This can be calculated iteratively as:
// P(k)_unique = 365/365 * 354/356 * 363/365 * ... (365-k)/365
// P(k)_collision = 1 - P(k)_unique
//
// Because of multiplying many small numbers,
// there is high chance of underflow that's why logarithm are used
// to change multiplication into addition
func birthdayProblemTheoretical(n, k int64) float64 {
	// res - holds value of p(k)_unique
	res := .0
	for i := 1; i < int(k); i++ {
		// We're dealing with logarithms
		// this subtraction represents division log(a/b) = log(a) - log(b)
		r := math.Log(float64(int(n)-i)) - math.Log(float64(n))
		// this addition represents multiplication log(a*b) = log(a) + log(b)
		res += r
	}

	// Let's refresh logarithm definition:
	//		log_b(a) = c    <==>   a = b^c
	//
	// We have computed value `c`, that is stored in `res` variable,
	// To compute it we use natiral logarithm function which base `b` is `e`
	// to retrieve value `a` from logarithm, we need rise `b` to power of `c`.
	// In our case that's simply invocation of a function math.Exp()
	//
	// Rest of computation is calculation of complement.
	return 1 - math.Exp(res)
}

// There is alternative way to compute theoretical value of probability of collision
// It's writen concisely in [John D. Cook blog post](https://www.johndcook.com/blog/2016/01/30/general-birthday-problem/)
// which I recommend you to read.
//
// Alternative method translate mathematical equation one-to-one
// with addition of mathematical brilliance to overcome floating point underflow :)
//
// Solution to birthday problem can be written as:
//
// 		P(k) = 1-(n!/(n-k)!/n^k)			[E1]
//
// You can see a lot of factorials, for even for "small numbers" like 365,
// factorial of this number is huge and when this number used to divide by another big number
// there is chance of underflow occurring, to solve this we can use use logarithms.
//
// Golang has function Lgamma, that calculates logarithm value of gamma function.
// Gamma function is defined as Γ=(n-1)! which is exactly like factorial but minus one.
// We use this function in calculation bellow to compute equation [E1] using log gamma.
func birthdayProblemTheoretical2(n, k float64) float64 {
	res, _ := math.Lgamma(n + 1)
	res2, _ := math.Lgamma(n - k + 1)
	res -= res2
	res -= k * math.Log(n)

	return 1 - math.Exp(res)
}

func TestBirthdayProblemMonteCarlo(t *testing.T) {
	simulation := birthdayProblemMonteCarlo(1000, 365, 23)
	theoretical := birthdayProblemTheoretical(365, 23)
	theoretical2 := birthdayProblemTheoretical2(365, 23)

	between(t, simulation, 0.48, 0.51, 0.01)
	between(t, theoretical, 0.48, 0.51, 0.01)
	between(t, theoretical2, 0.48, 0.51, 0.01)
}

func TestBirthdayProblemPlot(t *testing.T) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Birthday problem - probability of two people having birthday on the same day"
	p.X.Label.Text = "k - number of people"
	p.Y.Label.Text = "p(k) probability"

	days := 365
	sampleSize := 300

	var resultMC, resultT1 plotter.XYs
	for i := 0; i < 100; i++ {
		resultMC = append(resultMC, plotter.XY{
			X: float64(i),
			Y: birthdayProblemMonteCarlo(sampleSize, days, i),
		})
		resultT1 = append(resultT1, plotter.XY{
			X: float64(i),
			Y: birthdayProblemTheoretical(int64(days), int64(i)),
		})
	}

	err = plotutil.AddLinePoints(p,
		fmt.Sprintf("Simulation (sample=%d) ", sampleSize), resultMC,
		"Theoretical", resultT1,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := p.Save(18*vg.Inch, 9*vg.Inch, "birthday_problem_calo_test.png"); err != nil {
		t.Fatal(err)
	}
}
