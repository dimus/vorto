package manager

import "github.com/dimus/vorto/domain/entity"

type Manager struct{}

func NewManager() Manager {
	return Manager{}
}

func (m Manager) AutoSelect(cs *entity.CardStack, bin entity.BinType) []entity.Card {
	var res []entity.Card
	return res
}

func (m Manager) Analyse(card *entity.Card) {
}

func (m Manager) Move(card *entity.Card) {
}
