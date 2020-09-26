package usecase

import "github.com/dimus/vorto/domain/entity"

type Loader interface {
	Init() error
	Load(set string) (*entity.CardStack, error)
	Save(cards *entity.CardStack) error
}

type Manager interface {
	AutoSelect(cs *entity.CardStack, bin entity.BinType) []entity.Card
	Analyse(card entity.Card)
	Move(card entity.Card)
}

type Teacher interface {
	Train(entity.BinType)
	Ask(card entity.Card) int
}

type Scorer interface {
	Score(value, answer string) int
	// BadScore grades a score as very bad (-1), bad (0), good (1)
	BadScore(score int) int
}
