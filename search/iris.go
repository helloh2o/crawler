/*
@时间 : 2020/10/15 10:27
@作者 : Admin
@描述 : //todo请描述详细的代码功能
*/
package search

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"os"
)

var app = iris.New()

func RunWeb(addr string) {
	app.Logger().SetLevel("debug")
	app.HandleDir("/", "./static")
	app.RegisterView(iris.HTML("./static", ".html"))
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut},
		AllowedHeaders:   []string{"*"},
	}))
	app.AllowMethods(iris.MethodOptions)
	//// 同时写文件日志与控制台日志
	//app.Logger().SetOutput(io.MultiWriter(f, os.Stdout))
	//// or 使用下面这个
	//// 日志只生成到文件
	app.Logger().SetOutput(os.Stdout)
	//////
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
