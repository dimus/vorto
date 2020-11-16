package entity

import (
	"strings"
	"time"
)

type CardStack struct {
	StackType `json:"stackType"`
	Set       string `json:"set"`
	Bins      `json:"bins"`
}

type Card struct {
	ID      string  `json:"id"`
	Val     string  `json:"value"`
	Def     string  `json:"defenition"`
	SortVal float32 `json:"-"`
	Reply   `json:"replies"`
}

type Reply struct {
	Answers   []bool `json:"answers"`
	TimeStamp int32  `json:"ts"`
}

func (rs Reply) Add(r bool) Reply {
	var res = []bool{r}
	res = append(res, rs.Answers...)
	if len(res) > 5 {
		res = res[0:5]
	}
	return Reply{Answers: res, TimeStamp: int32(time.Now().Unix())}
}

func (rs Reply) LastGoodAnsw() int {
	res := 0
	for _, v := range rs.Answers {
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
