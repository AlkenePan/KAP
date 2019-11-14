package web

import (
	"encoding/json"
	"github.com/kataras/iris"
	app2 "youzoo/why/pkg/app"
)

func StartApi(port string) {
	web := iris.Default()
	web.Logger().SetLevel("debug")

	web.Handle("GET", "/", func(ctx iris.Context) {
		asd, _ := json.Marshal(app2.Executable{"test"})
		response := string(asd)
		ctx.WriteString(response)
	})
	app_group := web.Party("/app")
	{
		app_group.Handle("POST", "/new", AppNew)
		app_group.Handle("POST", "/update", AppUpdate)
		app_group.Handle("GET", "/{appid}", AppFind)
		app_group.Handle("GET", "/list")
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