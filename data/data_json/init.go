package data_json

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dimus/vorto/domain/entity"
	"github.com/gnames/gnlib/sys"
)

func (e EngineJSON) Init() error {
	err := e.touchWorkDirs()
	if err != nil {
		return err
	}

	err = e.touchFiles()
	if err != nil {
		return err
	}

	err = e.initJSON()
	if err != nil {
		return err
	}

	return nil
}

func (e EngineJSON) touchWorkDirs() error {
	dir := e.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")
	for set := range e.Config.Sets {
		setDir := filepath.Join(cardsDir, set)
		err := sys.MakeDir(setDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e EngineJSON) touchFiles() error {
	dir := e.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")

	for set := range e.Config.Sets {
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
		err := e.touchCardStackFile(setDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e EngineJSON) touchCardStackFile(dir string) error {
	filePath := filepath.Join(dir, "card-stack.txt")
	if sys.FileExists(filePath) {
		return nil
	}

	text := "type = general\n"
	err := ioutil.WriteFile(filePath, []byte(text), 0644)
	return err
}

func (e EngineJSON) initJSON() error {
	dir := e.Config.DataDir
	cardsDir := filepath.Join(dir, "flashcards")

	for set := range e.Config.Sets {
		setDir := filepath.Join(cardsDir, set)
		filePath := filepath.Join(setDir, e.FileJSON)
		if sys.FileExists(filePath) {
			return nil
		}
		var emptyCardMap cardMap = make(map[string]entity.Replies)
		csJSON, err := e.Encoder.Encode(emptyCardMap)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filePath, csJSON, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
