package entity

type BinType int

const (
	New BinType = iota
	Learning
	Vocabulary
)

func (ft BinType) String() string {
	switch ft {
	case Vocabulary:
		return "vocabulary"
	case Learning:
		return "learning"
	case New:
		return "new"
	}
	return "n/a"
}
