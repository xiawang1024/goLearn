package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/rs/cors"
	"log"
	"time"
	"xorm.io/xorm"
)

type User struct {
	Id int64
	Name string `form:"name"`
	Age int `form:"age"`
	Sex string `xorm:"varchar(10)" form:"sex"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	CityInfo `xorm:"extends"`
}

type CityInfo struct {
	City string `form:"city"`
	Street string `form:"street"`
}

var engine *xorm.Engine
var err error

func main() {
	app:=iris.New()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	app.WrapRouter(c.ServeHTTP)

	app.UseRouter(logger.New())
	app.UseRouter(recover.New())
	engine,err=xorm.NewEngine("mysql","root:password@tcp(localhost:3306)/iris?charset?utf8")

	if err != nil {
		log.Fatalf("mysql database connections failed: %v",err)
	}

	err=engine.Sync2(new(User))
	if err != nil {
		log.Fatalf("user table sync failed: %v",err)
	}

	user:=app.Party("/user")
	{
		user.Get("/{name}", getUser)

		user.Post("/", addUser)
	}

	err=app.Listen(":8081")
	if err != nil {
		log.Fatalf("app listening on port 8080 is failed: %v",err)
	}
}

func addUser(ctx iris.Context) {
	var user User
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

	ctx.JSON(iris.Map{
		"code":0,
		"message":"success",
	})
}

func getUser(ctx iris.Context)  {
	name:=ctx.Params().GetDefault("name","王二狗")
	ctx.JSON(iris.Map{
		"code":0,
		"data":name,
	})
}


