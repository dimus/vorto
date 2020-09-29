package entity

import (
	"strings"
)

type CardStack struct {
	StackType `json:"stackType"`
	Set       string `json:"set"`
	Bins      `json:"bins"`
}

type Card struct {
	ID      string `json:"id"`
	Val     string `json:"value"`
	Def     string `json:"defenition"`
	Replies `json:"replies"`
}

type Replies []bool

func (rs Replies) Add(r bool) Replies {
	var res Replies = []bool{r}
	res = append(res, rs...)
	if len(res) <= 5 {
		return res
	}
	return res[0:5]
}

func (rs Replies) LastGoodAnsw() int {
	res := 0
	for _, v := range rs {
		if !v {
			break
		}
		res++
	}
	return res
}

type Bins map[BinType][]*Card

type StackType int

const (
	General StackType = iota
	Esperanto
)

func (st StackType) ProcessInput(input string) string {
	switch st {
	case Esperanto:
		return normalizeXNotation(input)
	default:
		return input
	}
}

func normalizeXNotation(input string) string {
	r := strings.NewReplacer("cx", "ĉ", "jx", "ĵ", "ux", "ŭ", "gx", "ĝ", "sx", "ŝ")
	return r.Replace(input)
}
