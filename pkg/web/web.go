package web

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"youzoo/why/pkg/storage"
)
var db, err = storage.OpenDb("/tmp/test.db")

func StartApi(port string) {
	web := iris.New()
	web.Logger().SetLevel("debug")
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},   //允许通过的主机名称
		AllowCredentials: true,
	})
	appGroup := web.Party("/app", crs).AllowMethods(iris.MethodOptions)
	{
		appGroup.Handle("POST", "/new", AppNew)
		appGroup.Handle("POST", "/update", AppUpdate)
		appGroup.Handle("GET", "/{appid}", AppFind)
		appGroup.Handle("GET", "/list", AppList)
	}
	keyGroup := web.Party("/key", crs).AllowMethods(iris.MethodOptions)
	{
		keyGroup.Handle("POST", "/update", KeyUpdate)
		keyGroup.Handle("GET", "/{appid}", KeyFind)
		keyGroup.Handle("GET", "/pub/{appid}", KeyFindPub)
		keyGroup.Handle("GET", "/pri/{appid}", KeyFindPri)

	}
	alertGroup := web.Party("/alert", crs).AllowMethods(iris.MethodOptions)
	{
		alertGroup.Handle("POST", "/new", AlertNew)
		alertGroup.Handle("POST", "/update", AlertUpdate)
		//alert_group.Handle("GET", "/search", AlertSearch)
		alertGroup.Handle("GET", "/list", AlertList)
	}
	buildGroup := web.Party("/build", crs).AllowMethods(iris.MethodOptions)
	{
		buildGroup.Handle("POST", "/new", BuildNew)
		buildGroup.Handle("POST", "/status/set", BuildStatusSet)
		buildGroup.Handle("POST", "/status/get/{buildID:int}", BuildStatusGet)
	}
	web.Run(iris.Addr(port), iris.WithoutServerError(iris.ErrServerClosed))

}


func ErrorHandling(err error, ctx iris.Context) bool {
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return true
	}
	return false
}