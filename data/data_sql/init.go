package data_sql

import (
	"fmt"
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
	}
	return nil
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
        (id TEXT PRIMARY KEY,
        value TEXT,
        description TEXT,
        bin INTEGER);
      CREATE INDEX IF NOT EXISTS cards_bin_idx
        ON cards (bin);

      DELETE FROM cards;

      CREATE TABLE IF NOT EXISTS stats
        (
          id TEXT NOT NULL,
          created_at TEXT NOT NULL,
          success INTEGER NOT NULL,
          PRIMARY KEY (id, created_at)
        );
      CREATE INDEX IF NOT EXISTS stats_created_at_idx
        ON stats (created_at);`
		_, err = esql.DB.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}
