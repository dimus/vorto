package data_sql

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dimus/vorto/domain/entity"
	"github.com/gnames/gnames/lib/gnuuid"
)

func (esql EngineSQL) Load(set string) (*entity.CardStack, error) {
	res, err := esql.initCardStack(set)
	if err != nil {
		return &res, err
	}

	err = esql.populateCardStack(&res)
	if err != nil {
		return &res, err
	}

	err = esql.populateDB(set, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

func (esql EngineSQL) populateCardStack(cs *entity.CardStack) error {
	dir := esql.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")
	setDir := filepath.Join(cardsDir, cs.Set)

	for _, file := range cardBins {
		err := esql.loadFile(setDir, file, cs)
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

func (esql EngineSQL) initCardStack(set string) (entity.CardStack, error) {
	res := entity.CardStack{
		Bins: make(map[entity.BinType][]entity.Card),
	}
	if _, ok := esql.Sets[set]; !ok {
		return res, fmt.Errorf("set '%s' is not in Sets", set)
	}
	res.Set = set
	return res, nil
}

func (esql EngineSQL) populateDB(set string, cs *entity.CardStack) error {
	var err error
	dir := esql.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")
	setDir := filepath.Join(cardsDir, set)
	fileDBPath := filepath.Join(setDir, esql.FileDB)
	esql.DB, err = sql.Open("sqlite3", fileDBPath)
	if err != nil {
		return err
	}
	tx, err := esql.DB.Begin()
	if err != nil {
		return err
	}
	count := 0
	stmt, err := tx.Prepare(`
	  INSERT OR IGNORE INTO cards (id, value, description, bin)
	  VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for bin, cards := range cs.Bins {
		for _, card := range cards {
			count++
			_, err = stmt.Exec(card.ID, card.Val, card.Def, int(bin))
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	if count == 0 {
		fmt.Printf("TODO: Add instruction how to load names\n")
	}
	return nil
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

func (esql EngineSQL) loadFile(dir string, bin entity.BinType, cs *entity.CardStack) error {
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
		cs.Bins[bin] = append(cs.Bins[bin], entity.Card{ID: gnuuid.New(val).String(), Val: val, Def: def})
	}
	return nil
}
