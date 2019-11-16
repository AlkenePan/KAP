package web

import (
	"github.com/kataras/iris/v12"
	app2 "youzoo/why/pkg/app"
	"youzoo/why/pkg/storage"
)

// POST /key/update
func KeyUpdate(ctx iris.Context) {
	var app app2.App
	err := ctx.ReadJSON(&app)
	if ErrorHandling(err, ctx) {
		return
	}
	// update CryptoTable
	err = storage.UpdateApp(app, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.WriteString(app.Dumps())
}

// GET /key/{appid}
func KeyFind(ctx iris.Context) {
	appid := ctx.Params().Get("appid")
	// find CryptoTable
	cryptoTable, err := storage.FindKeyPair(appid, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(cryptoTable)
}

// GET /key/pub/{appid}
func KeyFindPub(ctx iris.Context) {
	appid := ctx.Params().Get("appid")
	// find CryptoTable
	cryptoTable, err := storage.FindPubKey(appid, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(cryptoTable)
}

// GET /app/pri/{appid}
func KeyFindPri(ctx iris.Context) {
	appid := ctx.Params().Get("appid")
	// find CryptoTable
	cryptoTable, err := storage.FindPriKey(appid, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(cryptoTable)
}
