package data_json

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dimus/vorto/domain/entity"
	"github.com/fatih/color"
)

func (e EngineJSON) Save(cs *entity.CardStack) error {
	toLearn, vocab := partitionVocablulary(cs)
	toVocab, learn := partitionLearning(cs)
	cs.Bins[entity.Vocabulary] = append(vocab, toVocab...)
	cs.Bins[entity.Learning] = append(learn, toLearn...)

	newWordsNum := 25 - len(cs.Bins[entity.Learning])
	addNewWords(cs, newWordsNum)

	e.writeToFiles(cs)

	color.Green("\n\nVocabulary has %d words.\n", len(cs.Bins[entity.Vocabulary]))
	return nil
}

func addNewWords(cs *entity.CardStack, num int) {
	if num <= 15 {
		return
	}

	if num > len(cs.Bins[entity.New]) {
		num = len(cs.Bins[entity.New])
	}

	shuffleCards(cs.Bins[entity.New])
	newLearn := cs.Bins[entity.New][0:num]
	color.Green("\nMoving from 'new' to 'learning':")
	for i, v := range newLearn {
		color.Green("%d. %s", i, v.Val)
	}
	cs.Bins[entity.Learning] = append(cs.Bins[entity.Learning],
		newLearn...)
	cs.Bins[entity.New] = cs.Bins[entity.New][num:]
}

func shuffleCards(cards []*entity.Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func partitionVocablulary(cs *entity.CardStack) ([]*entity.Card, []*entity.Card) {
	var toLearn, vocab []*entity.Card
	for _, c := range cs.Bins[entity.Vocabulary] {
		if len(c.Reply.Answers) > 0 && !c.Reply.Answers[0] {
			var reply entity.Reply
			c.Reply = reply
			toLearn = append(toLearn, c)
		} else {
			vocab = append(vocab, c)
		}
	}
	if len(toLearn) > 0 {
		color.Red("\nMoving from 'vocabulary' to 'learning':")
	}
	for i, v := range toLearn {
		color.Red("%d. %s\n", i+1, v.Val)
	}
	return toLearn, vocab
}

func partitionLearning(cs *entity.CardStack) ([]*entity.Card, []*entity.Card) {
	var toVocab, learn []*entity.Card
	for _, c := range cs.Bins[entity.Learning] {
		if c.Reply.LastGoodAnsw() > 2 {
			toVocab = append(toVocab, c)
		} else {
			learn = append(learn, c)
		}
	}
	if len(toVocab) > 0 {
		color.Green("\nMoving from learning to vocabulary:")
	}
	for i, v := range toVocab {
		color.Green("%d. %s\n", i+1, v.Val)
	}
	return toVocab, learn
}

func (e EngineJSON) writeToFiles(cs *entity.CardStack) error {
	var cMap cardMap = make(map[string]entity.Reply)
	bins := []entity.BinType{entity.Vocabulary, entity.Learning, entity.New}
	for _, bin := range bins {
		err := e.saveBin(cs, bin, cMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e EngineJSON) saveBin(cs *entity.CardStack, bin entity.BinType, cMap cardMap) error {
	path := filepath.Join(e.DataDir, "flashcards", cs.Set, bin.String()+".txt")
	var cards []string
	for _, card := range cs.Bins[bin] {
		if _, ok := cMap[card.Val]; ok {
			fmt.Printf("Card '%s' already exists, ignoring.\n", card.Val)
			continue
		}
		cMap[card.Val] = card.Reply
		cards = append(cards, fmt.Sprintf("%s = %s", card.Val, card.Def))
	}

	sort.Strings(cards)
	err := ioutil.WriteFile(path, []byte(strings.Join(cards, "\n")), 0664)
	if err != nil {
		return err
	}
	return e.saveCardMap(cs, cMap)
}

func (e EngineJSON) saveCardMap(cs *entity.CardStack, m cardMap) error {
	path := filepath.Join(e.DataDir, "flashcards", cs.Set, e.FileJSON)
	res, err := e.Encoder.Encode(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, res, 0644)
}
