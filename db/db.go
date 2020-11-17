package db

import (
	"log"
	"myapp/models"
	"xorm.io/xorm"
)

var Engine *xorm.Engine
var err error


func init()  {
	Engine,err=xorm.NewEngine("mysql","root:password@tcp(localhost:3306)/iris?charset?utf8")

	if err != nil {
		log.Fatalf("mysql database connections failed: %v",err)
	}

	err=Engine.Sync2(new(models.User))
	if err != nil {
		log.Fatalf("user table sync failed: %v",err)
	}
}