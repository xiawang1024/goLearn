package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"myapp/models"
	"myapp/db"
)

type validationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
}

type User = models.User

var err error



func AddUser(ctx iris.Context) {
	var user User
	session := sessions.Get(ctx)

	err:=ctx.ReadBody(&user)

	ctx.Application().Logger().Infof("%v",user)

	if err != nil {


		if errs,ok:=err.(validator.ValidationErrors);ok {
			validationErrors := wrapValidationErrors(errs)
			ctx.StopWithProblem(iris.StatusBadRequest,iris.NewProblem().Title("validation error").Detail("one or more fields failed to be validated").Type("/user/validation-errors").Key("errors",validationErrors))
			return
		}
		ctx.StopWithError(iris.StatusBadRequest,err)


		ctx.Application().Logger().Errorf("v%",err)
		return
	}

	ctx.Application().Logger().Infof("User: %#+v",user)

	_,err=db.Engine.Insert(&user)

	if err != nil {
		ctx.Application().Logger().Errorf("user table insert failed: %v",err)
		ctx.StopWithError(iris.StatusBadRequest,err)
	}
	session.Set("user",&user)
	ctx.JSON(iris.Map{
		"code":0,
		"message":"success",
	})
}

func GetUser(ctx iris.Context)  {


	session := sessions.Get(ctx)
	//多条记录查询
	name:=ctx.Params().Get("name")
	users := make([]User,0)
	err=db.Engine.Where("name=?",name).Find(&users)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest,err)
		ctx.Application().Logger().Errorf("user find by name failed: %v",err)
	}

	ctx.JSON(iris.Map{
		"code":0,
		"data":users,
		"session":session.Get("user"),
	})
}

func LoginOut(ctx iris.Context) {
	session := sessions.Get(ctx)
	session.Destroy()

	ctx.JSON(iris.Map{
		"code":0,
		"msg":"success",
	})

}



func wrapValidationErrors(errs validator.ValidationErrors) []validationError {
	validationErrors := make([]validationError, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, validationError{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
		})
	}

	return validationErrors
}
