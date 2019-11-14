package app

import (
	"github.com/google/uuid"
	"testing"
)

func TestApp(t *testing.T) {
	execuable := Executable{"~/.bin/test"}
	source := Source{"go"}
	app := App{uuid.New(), execuable, source}
	appDumped := app.Dumps()
	appLoaded := new(App)
	appLoaded.Loads(appDumped)
	if appLoaded.Appid != app.Appid {
		t.Error(
			"expected", app.Appid,
			"got", appLoaded.Appid,
		)
	}

}
