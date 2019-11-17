package web

import (
	"github.com/google/uuid"
	"github.com/kataras/iris"
	app2 "youzoo/why/pkg/app"
	"youzoo/why/pkg/crypto"
	"youzoo/why/pkg/storage"
)

// POST /app/new
func AppNew(ctx iris.Context) {
	var app app2.App
	err := ctx.ReadJSON(&app)
	app.Appid = uuid.New()
	if ErrorHandling(err, ctx) {
		return
	}
	// create AppTable
	err = storage.CreateApp(app, db)
	pri, pub := crypto.GenerateKeyPair(2048)
	_ = storage.NewKeyPair(app, string(crypto.PublicKeyToBytes(pub)[:]), string(crypto.PrivateKeyToBytes(pri)[:]), db)
	if ErrorHandling(err, ctx) {
		return
	}
	_, _ = ctx.WriteString(app.Dumps())
	return
}

// POST /app/update
func AppUpdate(ctx iris.Context) {
	var app app2.App
	err := ctx.ReadJSON(&app)
	if ErrorHandling(err, ctx) {
		return
	}
	// update AppTable
	err = storage.UpdateApp(app, db)
	if ErrorHandling(err, ctx) {
		return
	}
	_, _ = ctx.WriteString(app.Dumps())
}

// GET /app/{appid}
func AppFind(ctx iris.Context) {
	appid := ctx.Params().Get("appid")
	// find AppTable
	app, err := storage.FindApp(appid, db)
	if ErrorHandling(err, ctx) {
		return
	}
	_, _ = ctx.JSON(app)
}

// GET /app/list?from=<int>&count=<int>
func AppList(ctx iris.Context) {
	from := ctx.URLParamIntDefault("from", 0)
	count := ctx.URLParamIntDefault("count", 50)
	// search AppTable
	apps, err := storage.ListApp(from, count, db)
	if ErrorHandling(err, ctx) {
		return
	}
	_, _ = ctx.JSON(apps)
}
