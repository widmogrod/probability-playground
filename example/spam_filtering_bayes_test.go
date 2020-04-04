package example

import (
	"strings"
	"testing"
)

func TestSpamFilteringBayesTest(t *testing.T) {
	type train struct {
		spam bool
		text string
	}
	type test struct {
		text  string
		pSpam float64
		pHam  float64
	}

	useCases := map[string]struct {
		train []train
		test  test
	}{
		"toy": {
			train: []train{
				{true, "send me your password"},
				{false, "send me your picture"},
				{true, "what is your password"},
				{false, "what is your name"},
			},
			test: test{
				text:  "what is your password",
				pSpam: 0,
				pHam:  0,
			},
		},
	}
	for name, uc := range useCases {
		t.Run(name, func(t *testing.T) {
			// prepare
			bowc := BowClass{}
			for _, sample := range uc.train {
				words := strings.Split(sample.text, " ")
				for _, w := range words {
					bowc.Inc(sample.spam, w)
				}
			}

			// train
			//                   P(H) * P(E|H)
			// P(H|E) = --------------------------------
			//           P(H) * P(E|H) + P(-H) * P(E|-H)

			PrSpamInit := 0.5

			words := strings.Split(uc.test.text, " ")
			for _, w := range words {
				if !bowc.Has(w) {
					continue
				}

				PrSpamInit = (PrSpamInit * bowc.Proportion(true, w)) / (PrSpamInit*bowc.Proportion(true, w) + ((1 - PrSpamInit) * bowc.Proportion(false, w)))

				t.Logf("                  %f * %f", PrSpamInit, bowc.Proportion(true, w))
				t.Logf("P(H|E) = --------------------------------")
				t.Logf("          %f * %f + %f * %f", PrSpamInit, bowc.Proportion(true, w), (1 - PrSpamInit), bowc.Proportion(false, w))
				t.Logf("")
				t.Logf("word: %s   P(spam)=%f\n\n", w, PrSpamInit)
			}

			PrHamInit := 0.5

			words = strings.Split(uc.test.text, " ")
			for _, w := range words {
				if !bowc.Has(w) {
					continue
				}

				PrHamInit = (PrHamInit * bowc.Proportion(false, w)) / (PrHamInit*bowc.Proportion(false, w) + ((1 - PrHamInit) * bowc.Proportion(true, w)))

				t.Logf("word: %s   P(ham)=%f", w, PrHamInit)
			}

			// test
		})
	}
}

type BoW map[string]float64

func (b BoW) Inc(w string) {
	if _, ok := b[w]; ok {
		b[w]++
	} else {
		b[w] = 1
	}
}

func (b BoW) Val(w string) float64 {
	if _, ok := b[w]; ok {
		return b[w]
	}

	return 0
}

func (b BoW) Total() float64 {
	result := .0

	// TODO: optimise!
	for _, count := range b {
		result += count
	}

	return result
}

type BowClass map[bool]BoW

func (bc BowClass) Inc(class bool, w string) {
	if b, ok := bc[class]; ok {
		b.Inc(w)
	} else {
		bc[class] = BoW{}
		bc[class].Inc(w)
	}
}

func (bc BowClass) Val(class bool, w string) float64 {
	if b, ok := bc[class]; ok {
		return b.Val(w)
	}

	return 0
}

func (bc BowClass) Has(w string) bool {
	if bc.Val(true, w) != .0 {
		return true
	}
	if bc.Val(false, w) != .0 {
		return true
	}

	return false
}

func (bc BowClass) Total(class bool) float64 {
	if b, ok := bc[class]; ok {
		return b.Total()
	}

	return 0
}

func (bc BowClass) Proportion(class bool, w string) float64 {
	return bc.Val(class, w) / bc.Total(class)
}
