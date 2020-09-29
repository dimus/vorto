package vorto

import (
	"github.com/dimus/vorto/config"
	"github.com/dimus/vorto/data/data_json"
	"github.com/dimus/vorto/domain/entity"
	"github.com/dimus/vorto/domain/usecase"
	"github.com/dimus/vorto/manager"
	"github.com/dimus/vorto/teacher"
	"github.com/fatih/color"
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
		Loader:  data_json.NewEngineJSON(cfg),
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
	color.Yellow("Learning new terms...")
	t.Train(entity.Learning)

	if len(cs.Bins[entity.Learning]) >= 15 {
		return
	}

	color.Green("\nChecking learned before words...\n")
	t.Train(entity.Vocabulary)
}

func (vrt Vorto) Save(cs *entity.CardStack) error {
	return vrt.Loader.Save(cs)
}
