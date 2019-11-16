package web

import (
	"github.com/kataras/iris/v12"
	"youzoo/why/pkg/storage"
)

// POST /build/new
func BuildNew(ctx iris.Context) {
	var buildTable storage.BuildTable
	err := ctx.ReadJSON(&buildTable)
	if ErrorHandling(err, ctx) {
		return
	}
	// create AppTable
	build, err := storage.CreateBuild(buildTable, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(build)
	return
}

// POST /build/status/set
func BuildStatusSet(ctx iris.Context) {
	var buildStatusJson storage.BuildStatusJson
	err := ctx.ReadJSON(&buildStatusJson)
	if ErrorHandling(err, ctx) {
		return
	}
	// update AppTable
	buildTable, err := storage.SetBuildStatus(buildStatusJson, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(buildTable)
}

// GET /build/status/get/{buildID}
func BuildStatusGet(ctx iris.Context) {
	buildId := ctx.Params().GetIntDefault("buildID", 0)
	if ErrorHandling(err, ctx) {
		return
	}
	// update AppTable
	buildTable, err := storage.GetBuildStatus(buildId, db)
	if ErrorHandling(err, ctx) {
		return
	}
	ctx.JSON(buildTable)
}
