package app

import (
	"encoding/json"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"log"
)

type App struct {
	Appid uuid.UUID
	ExecInfo Executable
	SourceInfo Source
}

type Executable struct {
	AbsPath string
}

type Source struct {
	Language string
}

// dumps to json
func (app App) Dumps() string {
	b, err := json.Marshal(app)
	if err != nil {
		log.Fatalln(err)
	}
	return string(b)
}

// loads from string
func (app *App) Loads(appJson string) {
	err := json.Unmarshal([]byte(appJson), &app)
	if err != nil {
		log.Fatalln(err)
	}
}
