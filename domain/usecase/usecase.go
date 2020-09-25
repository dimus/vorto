package usecase

import "github.com/dimus/vorto/domain/entity"

type Loader interface {
	Init() error
	Load(set string) (*entity.CardStack, error)
	Save(cards *entity.CardStack) error
}

type Manager interface {
	Analyse(card entity.Card)
	Move(card entity.Card)
}

type Teacher interface {
	Ask(card entity.Card, answer *string)
	Grade(card entity.Card)
}
