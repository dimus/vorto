package data_sql

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dimus/vorto/domain/entity"
	"github.com/gnames/gnames/lib/gnuuid"
)

func (esql EngineSQL) Load(set string) (*entity.CardStack, error) {
	res := entity.CardStack{
		Bins: make(map[entity.BinType][]entity.Card),
	}
	if _, ok := esql.Sets[set]; !ok {
		return &res, fmt.Errorf("set '%s' is not in Sets", set)
	}

	dir := esql.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")
	setDir := filepath.Join(cardsDir, set)

	for _, file := range cardBins {
		err := esql.loadFile(setDir, file, &res)
		if err != nil {
			return &res, err
		}
	}
	err := esql.populateDB(set, &res)
	if err != nil {
		return &res, err
	}
	return &res, nil
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
	  INSERT INTO cards (id, value, description, bin)
	  VALUES (?, ?, ?, ?)`)
	defer stmt.Close()

	for bin, cards := range cs.Bins {
		for _, card := range cards {
			count++
			stmt.Exec(card.ID, card.Val, card.Def, int(bin))
		}
	}
	tx.Commit()
	if count == 0 {
		fmt.Printf("TODO: Add instruction how to load names\n")
	}
	return nil
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
		fields := strings.Split(line, " = ")
		if len(fields) != 2 {
			continue
		}
		val := strings.Trim(fields[0], " ")
		def := strings.Trim(fields[1], " ")
		cs.Bins[bin] = append(cs.Bins[bin], entity.Card{ID: gnuuid.New(val).String(), Val: val, Def: def})
	}
	return nil
}
