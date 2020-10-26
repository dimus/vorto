package teacher

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dimus/vorto/domain/entity"
	"github.com/dimus/vorto/domain/usecase"
	"github.com/dimus/vorto/score"
	"github.com/fatih/color"
)

type Teacher struct {
	*entity.CardStack
	usecase.Scorer
	Score int
}

func NewTeacher(cs *entity.CardStack) Teacher {
	return Teacher{
		CardStack: cs,
		Scorer:    score.ScorePunish{},
	}
}

func (t Teacher) Train(bin entity.BinType) {
	var ok bool
	var cards []*entity.Card
	if cards, ok = t.CardStack.Bins[bin]; ok && len(cards) > 0 {
		cards = selectCards(cards)
		t.processCards(cards)
	} else {
		log.Printf("There are no cards in a '%s' bin.", bin)
	}
}

func selectCards(cards []*entity.Card) []*entity.Card {
	if len(cards) < 25 {
		shuffleCards(cards)
		return cards
	}
	bad, good, perfect := partitionCards(cards)
	return examCards(bad, good, perfect)
}

func (t Teacher) Ask(card *entity.Card) int {
	scoreFinal := 0
	fmt.Printf("What is: %s\n", card.Def)
	for {
		fmt.Print("-> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = t.ProcessInput(text)
		score := t.Scorer.Score(text, card.Val)
		if t.BadScore(score) == -1 {
			scoreFinal += score
			color.Red("!> %s", card.Val)
			color.Red("Score %d", score)
		} else if t.BadScore(score) == 0 {
			scoreFinal += score
			color.Yellow("~> %s", card.Val)
			color.Yellow("Score %d", score)
		} else {
			color.Green("-> %s", card.Val)
			color.Green("Correct!")
			break
		}
	}
	// If BadScore returns 1, the answer is correct.
	card.Replies = card.Replies.Add(t.BadScore(scoreFinal) == 1)
	return scoreFinal
}

func (t Teacher) processCards(cards []*entity.Card) {
	total := 0
	shuffleCards(cards)
	for i, card := range cards {
		fmt.Printf("\n%d out of %d\n", i+1, len(cards))
		score := t.Ask(card)
		total += score
		fmt.Printf("Total so far: %d\n", total)
		fmt.Printf("res: %+v\n", card.Replies)
	}
}

func examCards(bad, good, perfect []*entity.Card) []*entity.Card {
	var res []*entity.Card
	var leftowers []*entity.Card
	if len(perfect) < 7 {
		res = append(res, perfect...)
	} else {
		res = append(res, perfect[0:7]...)
		leftowers = append(leftowers, perfect[7:]...)
	}
	if len(good) < 7 {
		res = append(res, good...)
	} else {
		res = append(res, good[0:7]...)
		leftowers = append(leftowers, good[7:]...)
	}
	badSize := 25 - len(res)
	if len(bad) < badSize {
		res = append(res, bad...)
	} else {
		res = append(res, bad[0:badSize]...)
	}
	if len(res) > 24 || len(leftowers) == 0 {
		shuffleCards(res)
		return res
	}
	fillSize := 25 - len(res)
	if len(leftowers) < fillSize {
		res = append(res, leftowers...)
	} else {
		res = append(res, leftowers[0:fillSize]...)
	}
	shuffleCards(res)
	return res
}

func partitionCards(cards []*entity.Card) ([]*entity.Card, []*entity.Card, []*entity.Card) {
	var bad, good, perfect []*entity.Card
	for _, v := range cards {
		card := v
		switch goodAnswers(card) {
		case 0, 1, 2, 3:
			bad = append(bad, card)
		case 4:
			good = append(good, card)
		default:
			perfect = append(perfect, card)
		}
	}
	shuffleCards(bad)
	shuffleCards(good)
	shuffleCards(perfect)
	return bad, good, perfect
}

func goodAnswers(card *entity.Card) int {
	count := 0
	for _, v := range card.Replies {
		if !v {
			break
		}
		count += 1
	}
	return count
}

func shuffleCards(cards []*entity.Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}
