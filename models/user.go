package models

import "time"

type User struct {
	Id int64 `json:"id"`
	Name string `json:"name" validate:"required"`
	Age int `json:"age" validate:"gte=0,lte=120"`
	Sex string `xorm:"varchar(10)" json:"sex"`
	Created time.Time `xorm:"created" json:"-"`
	Updated time.Time `xorm:"updated" json:"-"`
	CityInfo `xorm:"extends"`
}

type CityInfo struct {
	City string `json:"city" validate:"required"`
	Street string `json:"street"`
}
