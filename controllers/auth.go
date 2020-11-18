package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"myapp/middleware"
)

func LogIn(ctx iris.Context)  {
	standardClaims := jwt.Claims{Issuer: "xiawang1024.com",Audience:jwt.Audience{"xiwang1024","gyy"}}
	customClaims := middleware.UserClaims{
		Claims:   middleware.J.Expiry(standardClaims),
		Username: "wangxia",
	}

	//j.WriteToken(ctx, customClaims)
	token ,_:= middleware.J.Token(customClaims)
	ctx.JSON(iris.Map{
		"code":0,
		"token":token,
	})
}
