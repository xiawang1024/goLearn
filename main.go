package main

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
	"github.com/rs/cors"
	"myapp/controllers"
	_ "myapp/db"
	"os"
	"regexp"
	"time"
)



type UserClaims struct {
	jwt.Claims
	Username string
}


var err error
var j *jwt.JWT

const cookieNameForSessionID = "session_id_cookie"

func main() {
	app:=iris.New()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	j=jwt.HMAC(15*time.Minute,"secret")

	app.WrapRouter(c.ServeHTTP)
	f,_:=os.Create("iris.log")
	app.Logger().SetOutput(f)

	app.UseRouter(logger.New())
	app.UseRouter(recover.New())
	sess := sessions.New(sessions.Config{Cookie: cookieNameForSessionID, AllowReclaim: true,Expires: 15 * time.Minute})
	app.Use(sess.Handler())
	app.Use(iris.Cache(15*time.Second))


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

	app.Validator = validator.New()

	user:=app.Party("/user",jwtMiddle)
	{
		user.Get("/{name}", controllers.GetUser)

		user.Post("/", controllers.AddUser)

		user.Get("/out",controllers.LoginOut)


	}

	app.Get("/auth", func(ctx iris.Context) {
		standardClaims := jwt.Claims{Issuer: "xiawang1024.com",Audience:jwt.Audience{"xiwang1024","gyy"}}
		customClaims := UserClaims{
			Claims:   j.Expiry(standardClaims),
			Username: "wangxia",
		}

		//j.WriteToken(ctx, customClaims)
		token ,_:= j.Token(customClaims)
		ctx.JSON(iris.Map{
			"code":0,
			"token":token,
		})
	})

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

func jwtMiddle(ctx iris.Context)  {
	var claims UserClaims
	err = j.VerifyToken(ctx,&claims)
	if err != nil {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		//ctx.StopWithError(iris.StatusUnauthorized,err)
	}
	ctx.Next()
}

