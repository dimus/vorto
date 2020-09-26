package data_sql

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"database/sql"

	"github.com/gnames/gnames/lib/sys"
	_ "github.com/mattn/go-sqlite3"
)

func (esql EngineSQL) Init() error {
	err := esql.touchWorkDirs()
	if err != nil {
		return err
	}

	err = esql.touchFiles()
	if err != nil {
		return err
	}

	err = esql.initDB()
	if err != nil {
		return err
	}

	return nil
}

func (esql EngineSQL) touchWorkDirs() error {
	dir := esql.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")
	for set := range esql.Config.Sets {
		setDir := filepath.Join(cardsDir, set)
		err := sys.MakeDir(setDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (esql EngineSQL) touchFiles() error {
	dir := esql.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")

	for set := range esql.Config.Sets {
		setDir := filepath.Join(cardsDir, set)
		for _, file := range cardBins {
			file := file.String() + ".txt"
			filePath := filepath.Join(setDir, file)
			if !sys.FileExists(filePath) {
				fmt.Printf("Creating file %s\n", filePath)
				f, err := os.Create(filePath)
				if err != nil {
					return err
				}
				f.Close()
			}
		}
		err := esql.touchCardStackFile(setDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (esql EngineSQL) touchCardStackFile(dir string) error {
	filePath := filepath.Join(dir, "card-stack.txt")
	if sys.FileExists(filePath) {
		return nil
	}

	text := "type = general\n"
	err := ioutil.WriteFile(filePath, []byte(text), 0644)
	return err
}

func (esql EngineSQL) initDB() error {
	var err error
	dir := esql.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")

	for set := range esql.Config.Sets {
		setDir := filepath.Join(cardsDir, set)
		fileDBPath := filepath.Join(setDir, esql.FileDB)
		esql.DB, err = sql.Open("sqlite3", fileDBPath)
		if err != nil {
			return err
		}

		q := `
      CREATE TABLE IF NOT EXISTS cards
        (id BLOB PRIMARY KEY,
        value TEXT,
        description TEXT,
        bin INTEGER);
      CREATE INDEX IF NOT EXISTS cards_bin_idx
        ON cards (bin);

      DELETE FROM cards;

      CREATE TABLE IF NOT EXISTS stats
        (
          id BLOB NOT NULL,
          success INTEGER NOT NULL,
          start_time INTEGER,
          start_typing_time INTEGER,
          end_typing_time INTEGER,
          PRIMARY KEY (id, start_time)
        );
      CREATE INDEX IF NOT EXISTS stats_start_time_idx
        ON stats (start_time);`
		_, err = esql.DB.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}
