package data_json

import (
	"github.com/dimus/vorto/config"
	"github.com/dimus/vorto/domain/entity"
	"github.com/gnames/gnlib/encode"
)

var cardBins = []entity.BinType{entity.Vocabulary, entity.Learning, entity.New}

type EngineJSON struct {
	config.Config
	FileJSON string
	Encoder  encode.GNjson
}

func NewEngineJSON(cfg config.Config) EngineJSON {
	return EngineJSON{
		Config:   cfg,
		FileJSON: "vorto.json",
		Encoder:  encode.GNjson{Pretty: true},
	}
}

// cardMap is a map of a card value to card replies.
type cardMap map[string]entity.Reply
