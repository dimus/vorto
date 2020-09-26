package vorto

import (
	"github.com/dimus/vorto/config"
	"github.com/dimus/vorto/data/data_sql"
	"github.com/dimus/vorto/domain/entity"
	"github.com/dimus/vorto/domain/usecase"
	"github.com/dimus/vorto/manager"
	"github.com/dimus/vorto/teacher"
)

type Vorto struct {
	config.Config
	usecase.Loader
	usecase.Teacher
	usecase.Manager
}

func NewVorto(cfg config.Config) Vorto {
	return Vorto{
		Config:  cfg,
		Loader:  data_sql.NewEngineSQL(cfg),
		Manager: manager.NewManager(),
	}
}

func (vrt Vorto) Init() error {
	return vrt.Loader.Init()
}

func (vrt Vorto) Load() (*entity.CardStack, error) {
	return vrt.Loader.Load(vrt.DefaultSet)
}

func (vrt Vorto) Run(cs *entity.CardStack) {
	t := teacher.NewTeacher(cs)
	vrt.Teacher = t
	t.Train(entity.Learning)
	t.Train(entity.Vocabulary)
	vrt.Save(cs)
}

func (vrt Vorto) Save(cs *entity.CardStack) error {
	return vrt.Loader.Save(cs)
}
