package entity

import (
	"strings"
	"time"
)

type CardStack struct {
	StackType
	Set string
	Bins
}

type Card struct {
	ID         string
	Val        string
	Def        string
	Replies    []Reply
	ReplyPause int
	ReplySpeed int
}

type Reply struct {
	Val bool
	time.Time
}

type Bins map[BinType][]Card

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
