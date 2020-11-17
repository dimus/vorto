package data_json

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dimus/vorto/domain/entity"
	"github.com/gnames/gnlib/gnuuid"
)

func (e EngineJSON) Load(set string) (*entity.CardStack, error) {
	res, err := e.initCardStack(set)
	if err != nil {
		return &res, err
	}

	err = e.populateCardStack(&res)
	if err != nil {
		return &res, err
	}
	return &res, nil
}

func (e EngineJSON) initCardStack(set string) (entity.CardStack, error) {
	res := entity.CardStack{
		Bins: make(map[entity.BinType][]*entity.Card),
	}
	if _, ok := e.Sets[set]; !ok {
		return res, fmt.Errorf("set '%s' is not in Sets", set)
	}
	res.Set = set
	return res, nil
}

func (e EngineJSON) populateCardStack(cs *entity.CardStack) error {
	dir := e.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")
	setDir := filepath.Join(cardsDir, cs.Set)

	for _, file := range cardBins {
		err := e.loadFile(setDir, file, cs)
		if err != nil {
			return err
		}
	}
	cs.StackType = stackType(setDir)
	return nil
}

func stackType(dir string) entity.StackType {
	filePath := filepath.Join(dir, "card-stack.txt")
	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Could not open file to get StackType: %s, setting it to 'General'.\n", err)
		return entity.General
	}
	lines := strings.Split(string(text), "\n")
	for _, line := range lines {
		key, val, success := parseLine(line)
		if !success || key != "type" {
			continue
		}
		switch val {
		case "esperanto", "Esperanto":
			return entity.Esperanto
		default:
			return entity.General
		}
	}
	log.Println("Could not get StackType from file, setting it to 'General'")
	return entity.General
}

func (e EngineJSON) savedCardMap(set string) (cardMap, error) {
	var res []CardStorage
	var cm cardMap = make(map[string]entity.Reply)
	filePath := filepath.Join(e.DataDir, "flashcards", set, e.FileJSON)
	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		return cm, err
	}
	err = e.Encoder.Decode(text, &res)
	for _, v := range res {
		cm[v.Value] = v.Reply
	}
	return cm, err
}

func (e EngineJSON) prepareDataJSON() {
}

// parseLine takes a string and breaks it into two fields if it detects
// ' = ' pattern. It ignores comments designated by leading '#' character
// and empty lines. It returns 2 fields, and a boolean that indicates
// a success or failure of parsing.
func parseLine(line string) (string, string, bool) {
	fields := strings.Split(line, " = ")
	if len(fields) != 2 {
		return "", "", false
	}

	field1 := strings.Trim(fields[0], " ")
	field2 := strings.Trim(fields[1], " ")
	if len(field1) < 1 || len(field2) < 1 || field1[0] == '#' {
		return "", "", false
	}
	return field1, field2, true
}

func (e EngineJSON) loadFile(dir string, bin entity.BinType, cs *entity.CardStack) error {
	ts := int32(time.Now().Unix())
	oldCardMap, err := e.savedCardMap(cs.Set)
	if err != nil {
		return err
	}
	filePath := filepath.Join(dir, bin.String()+".txt")
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		val, def, success := parseLine(line)
		if !success {
			continue
		}

		card := entity.Card{ID: gnuuid.New(val).String(), Val: val, Def: def}
		if reply, ok := oldCardMap[card.Val]; ok {
			card.Reply = reply
		}
		if bin == entity.Vocabulary && len(card.Reply.Answers) == 0 {
			card.Reply.Answers = []bool{true, true, true, true, true}
			card.Reply.TimeStamp = ts
		}
		cs.Bins[bin] = append(cs.Bins[bin], &card)
	}
	return nil
}
