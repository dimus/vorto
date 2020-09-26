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
	var cards []entity.Card
	if cards, ok = t.CardStack.Bins[bin]; ok && len(cards) > 0 {
		t.processCards(cards)
	} else {
		log.Printf("There are no cards in a '%s' bin.", bin)
	}
}

func (t Teacher) Ask(card entity.Card) int {
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
	return scoreFinal
}

func (t Teacher) processCards(cards []entity.Card) {
	total := 0
	shuffleCards(cards)
	for i, card := range cards {
		fmt.Printf("\n%d out of %d\n", i+1, len(cards))
		score := t.Ask(card)
		total += score
		fmt.Printf("Total so far: %d\n", total)
	}
}

func shuffleCards(cards []entity.Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}
