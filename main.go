package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
	"github.com/rs/cors"
	"log"
	"os"
	"regexp"
	"time"
	"xorm.io/xorm"
)

type User struct {
	Id int64 `json:"id"`
	Name string `form:"name" json:"name"`
	Age int `form:"age" json:"age"`
	Sex string `xorm:"varchar(10)" form:"sex" json:"sex"`
	Created time.Time `xorm:"created" json:"-"`
	Updated time.Time `xorm:"updated" json:"-"`
	CityInfo `xorm:"extends"`
}

type CityInfo struct {
	City string `form:"city" json:"city"`
	Street string `form:"street" json:"street"`
}

type UserClaims struct {
	jwt.Claims
	Username string
}

var engine *xorm.Engine
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

	engine,err=xorm.NewEngine("mysql","root:password@tcp(localhost:3306)/iris?charset?utf8")

	if err != nil {
		log.Fatalf("mysql database connections failed: %v",err)
	}

	err=engine.Sync2(new(User))
	if err != nil {
		log.Fatalf("user table sync failed: %v",err)
	}

	user:=app.Party("/user",jwtMiddle)
	{
		user.Get("/{name}", getUser)

		user.Post("/", addUser)

		user.Get("/out",loginOut)


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


	//err=app.Listen(":8080")
	//if err != nil {
	//	log.Fatalf("app listening on port 8080 is failed: %v",err)
	//}

	app.Run(iris.TLS(":443","mycert.crt","mykey.key"))
}

func addUser(ctx iris.Context) {
	var user User
	session := sessions.Get(ctx)

	err:=ctx.ReadBody(&user)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest,err)
		log.Fatalf("addUser error: %v",err)
		//return
	}

	ctx.Application().Logger().Infof("User: %#+v",user)

    _,err=engine.Insert(&user)

	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest,err)
		log.Fatalf("user table insert failed: %v",err)
	}
	session.Set("user",&user)
	ctx.JSON(iris.Map{
		"code":0,
		"message":"success",
		"session":session,
	})
}

func getUser(ctx iris.Context)  {
	//单条记录查询
	//name:=ctx.Params().Get("name")
	//user:=User{}
	//has,err:=engine.Where("name=?",name).Desc("id").Get(&user)
	//if err != nil {
	//	ctx.StopWithError(iris.StatusBadRequest,err)
	//	log.Fatalf("user Table find by name: %v",err)
	//}
	//if has {
	//	ctx.JSON(iris.Map{
	//		"code":0,
	//		"data":user,
	//	})
	//}else {
	//	ctx.JSON(iris.Map{
	//		"code":0,
	//		"message":"not found user",
	//	})
	//}

	//var claims UserClaims
	//if err :=j.VerifyToken(ctx,&claims);err!=nil{
	//	ctx.StopWithStatus(iris.StatusUnauthorized)
	//	return
	//}

	session := sessions.Get(ctx)
	//多条记录查询
	name:=ctx.Params().Get("name")
	users := make([]User,0)
	err=engine.Where("name=?",name).Find(&users)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest,err)
		log.Fatalf("user find by name failed: %v",err)
	}

	ctx.JSON(iris.Map{
		"code":0,
		"data":users,
		"session":session.Get("user"),
	})
}

func loginOut(ctx iris.Context) {
	session := sessions.Get(ctx)
	session.Destroy()

	ctx.JSON(iris.Map{
		"code":0,
		"msg":"success",
	})

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

