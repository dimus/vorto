package teacher

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dimus/vorto/domain/entity"
	"github.com/dimus/vorto/domain/usecase"
	"github.com/dimus/vorto/score"
	"github.com/fatih/color"
)

// the constants determine how may words of each cagegory will get into
// vocabulary questions.
const (
	// badNum is the number of questions with up to 3 correct answers
	badNum = 10
	// goodNum is the number of questions with 4 correct answers
	goodNum = 10
	// perfectNum is the number of questions with 5 correct anwers
	perfectNum = 5
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

func (t Teacher) Train(bin entity.BinType, withSecondChance bool) {
	var ok bool
	var cards []*entity.Card
	if cards, ok = t.CardStack.Bins[bin]; ok && len(cards) > 0 {
		cards = selectCards(cards)
		t.runExam(cards, withSecondChance)
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

func (t Teacher) Ask(card *entity.Card, withSecondChance bool) int {
	scoreFinal := 0
	fmt.Printf("What is: %s\n", card.Def)
	fmt.Printf("(%d answers, %d days ago)\n", goodAnswers(card), days(card))
	for {
		fmt.Print("-> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = t.ProcessInput(text)
		score := t.Scorer.Score(text, card.Val)

		// for vocabulary give one more chance if answer was wrong
		if t.BadScore(score) < 1 && withSecondChance {
			color.Yellow("Nope, one more chance")
			return t.Ask(card, false)
		}

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
	card.Reply = card.Reply.Add(t.BadScore(scoreFinal) == 1)
	return scoreFinal
}

func (t Teacher) runExam(cards []*entity.Card, withSecondChance bool) {
	total := 0
	shuffleCards(cards)
	for i, card := range cards {
		fmt.Printf("\n%d out of %d\n", i+1, len(cards))
		score := t.Ask(card, withSecondChance)
		total += score
		fmt.Printf("Total so far: %d\n", total)
		fmt.Printf("res: %+v\n", card.Reply.Answers)
	}
}

func examCards(bad, good, perfect []*entity.Card) []*entity.Card {
	var res []*entity.Card
	var leftowers []*entity.Card
	if len(perfect) < perfectNum {
		res = append(res, perfect...)
	} else {
		res = append(res, perfect[0:perfectNum]...)
		leftowers = append(leftowers, perfect[perfectNum:]...)
	}
	for _, v := range perfect {
		ts := strconv.Itoa(int(v.Reply.TimeStamp))
		ts = ts[len(ts)-5 : len(ts)-1]
	}
	if len(good) < goodNum {
		res = append(res, good...)
	} else {
		res = append(res, good[0:goodNum]...)
		leftowers = append(leftowers, good[goodNum:]...)
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
	perfect = preparePerfect(perfect)
	fmt.Printf("bad %d, good %d, perfect %d\n", len(bad), len(good), len(perfect))
	return bad, good, perfect
}

func preparePerfect(p []*entity.Card) []*entity.Card {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	rand.Seed(time.Now().UnixNano())
	minTS, diffTS := minDiffTS(p)
	type cardCarma struct {
		card  *entity.Card
		carma float32
	}
	for _, v := range p {
		fRand := r.Float32()
		var fTS float32
		if diffTS > 0 {
			fTS = (float32(v.Reply.TimeStamp - minTS)) / float32(diffTS)
		}
		v.SortVal = fRand + fTS
	}
	sort.Slice(p, func(i, j int) bool {
		return p[i].SortVal < p[j].SortVal
	})
	return p
}

func minDiffTS(p []*entity.Card) (int32, int32) {
	var minTS, maxTS int32
	for _, v := range p {
		if v.Reply.TimeStamp == 0.0 {
			continue
		}
		if minTS == 0 {
			minTS = v.Reply.TimeStamp
		}
		if maxTS == 0 {
			maxTS = v.Reply.TimeStamp
		}
		if minTS > v.Reply.TimeStamp {
			minTS = v.Reply.TimeStamp
		}
		if maxTS < v.Reply.TimeStamp {
			maxTS = v.Reply.TimeStamp
		}
	}
	return minTS, maxTS - minTS
}

func goodAnswers(card *entity.Card) int {
	count := 0
	for _, v := range card.Reply.Answers {
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

func days(card *entity.Card) int64 {
	t := time.Now().Unix()
	var diff int64
	if card.TimeStamp > 0 {
		diff = t - int64(card.TimeStamp)
	}
	return diff / 86400
}
