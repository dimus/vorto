package data_sql

import (
	"database/sql"

	"github.com/dimus/vorto/config"
	"github.com/dimus/vorto/domain/entity"
)

var cardBins = []entity.BinType{entity.Vocabulary, entity.Learning, entity.New}

type EngineSQL struct {
	config.Config
	FileDB string
	DB     *sql.DB
}

func NewEngineSQL(cfg config.Config) EngineSQL {
	esql := EngineSQL{
		Config: cfg,
		FileDB: "vorto.db",
	}
	return esql
}
