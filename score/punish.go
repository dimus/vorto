package score

import "github.com/gnames/gnames/lib/fuzzy"

type ScorePunish struct{}

func (s ScorePunish) Score(val, answer string) int {
	ed := fuzzy.EditDistance(val, answer)
	if ed < 4 {
		return ed * -1
	}

	return -10
}

func (s ScorePunish) BadScore(score int) int {
	if score == -10 {
		return -1
	}

	if score < 0 {
		return 0
	}

	return 1
}
