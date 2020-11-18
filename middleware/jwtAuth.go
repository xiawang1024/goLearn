package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"time"
)

type UserClaims struct {
	jwt.Claims
	Username string
}

var J *jwt.JWT

func init()  {
	J=jwt.HMAC(15*time.Minute,"secret")
}

func JwtAuth(ctx iris.Context)  {
	var claims UserClaims
	err := J.VerifyToken(ctx,&claims)
	if err != nil {
		ctx.StopWithStatus(iris.StatusUnauthorized)
	}
	ctx.Next()
}
