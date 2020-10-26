package score

import "github.com/gnames/gnlib/fuzzy"

type ScorePunish struct{}

var failure = -100

func (s ScorePunish) Score(val, answer string) int {
	ed := fuzzy.EditDistance(val, answer)
	if ed < 4 {
		return ed * -1
	}

	return failure
}

func (s ScorePunish) BadScore(score int) int {
	if score == failure {
		return -1
	}

	if score < 0 {
		return 0
	}

	return 1
}
