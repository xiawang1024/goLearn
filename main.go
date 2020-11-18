package main

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
	"github.com/rs/cors"
	"myapp/controllers"
	_ "myapp/db"
	"myapp/middleware"
	"os"
	"regexp"
	"time"
)

const cookieNameForSessionID = "session_id_cookie"

func main() {
	app:=iris.New()

	//跨域设置
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	app.WrapRouter(c.ServeHTTP)

	//日志
	f,_:=os.Create("iris.log")
	app.Logger().SetOutput(f)
	app.UseRouter(logger.New())
	app.UseRouter(recover.New())

	//设置session
	sess := sessions.New(sessions.Config{Cookie: cookieNameForSessionID, AllowReclaim: true,Expires: 15 * time.Minute})
	app.Use(sess.Handler())
	app.Use(iris.Cache(60*time.Second))

	//挂载静态资源
	app.HandleDir("/public",iris.Dir("./public"),iris.DirOptions{
		IndexName: "index.html",
		PushTargetsRegexp: map[string]*regexp.Regexp{
			"/":iris.MatchCommonAssets,
		},
		Cache: iris.DirCacheOptions{
			Enable: true,
			Encodings: []string{"gzip"},
			CompressIgnore: iris.MatchImagesAssets,
			CompressMinSize: 30 * iris.B,
		},
	})


	//数据绑定及验证
	app.Validator = validator.New()

	app.Get("/auth", controllers.LogIn)

	user:=app.Party("/user",middleware.JwtAuth)
	{
		user.Get("/{name}", controllers.GetUser)

		user.Post("/", controllers.AddUser)

		user.Get("/out",controllers.LoginOut)


	}


	cacheTest := app.Party("/cache")
	{
		cacheTest.Get("/", func(ctx iris.Context) {
			ctx.Header("X-Custom", "my  custom header")
			ctx.Writef("Hello World! %s", time.Now())
		})
	}


	//err=app.Listen(":8080")
	//if err != nil {
	//	log.Fatalf("app listening on port 8080 is failed: %v",err)
	//}

	app.Run(iris.TLS(":443","mycert.crt","mykey.key"))
}



