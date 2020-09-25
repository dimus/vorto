package vorto

import (
	"github.com/dimus/vorto/config"
	"github.com/dimus/vorto/data/data_sql"
	"github.com/dimus/vorto/domain/entity"
	"github.com/dimus/vorto/domain/usecase"
)

type Vorto struct {
	config.Config
	usecase.Loader
}

func NewVorto(cfg config.Config) Vorto {
	return Vorto{
		Config: cfg,
		Loader: data_sql.NewEngineSQL(cfg),
	}
}

func (vrt Vorto) Init() error {
	return vrt.Loader.Init()
}

func (vrt Vorto) Load() (*entity.CardStack, error) {
	return vrt.Loader.Load(vrt.DefaultSet)
}

func (vrt Vorto) Save(cs *entity.CardStack) error {
	return vrt.Loader.Save(cs)
}

func (vrt Vorto) Run(cs *entity.CardStack) {
}
