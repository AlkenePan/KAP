package web

import (
	"github.com/kataras/iris"
	"youzoo/why/pkg/storage"
)



// POST /app/new
func AlertNew(ctx iris.Context) {
	var alertTable storage.AlertTable
	err := ctx.ReadJSON(&alertTable)
	if ErrorHandling(err, ctx) {
		return
	}
	// create AppTable
	alert, err := storage.CreateAlert(alertTable, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(alert)
	return
}

// POST /alert/update
func AlertUpdate(ctx iris.Context) {
	var alertTable storage.AlertTable
	err := ctx.ReadJSON(&alertTable)
	if ErrorHandling(err, ctx) {
		return
	}
	// update AppTable
	alert, err := storage.UpdateAlert(alertTable, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(alert)
}

// GET /alert/list?from=<int>&count=<int>
func AlertList(ctx iris.Context) {
	from := ctx.URLParamIntDefault("from", 0)
	count := ctx.URLParamIntDefault("count", 50)
	// search AppTable
	alerts, err := storage.ListAlert(from, count, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(alerts)
}
