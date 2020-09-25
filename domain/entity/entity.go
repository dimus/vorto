package entity

import (
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
