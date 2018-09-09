package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id         int
	Name       string
	User_order []*User_order `orm:"reverse(many)"`
}
type User_order struct {
	Id         int
	Order_data string `orm:size(100)`
	User       *User  `orm:"rel(fk)"`
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "fred:fred@tcp(172.16.0.10:3306)/gotest?charset=utf8")
	// 需要在init中注册定义的model
	orm.RegisterModel(new(User), new(User_order))
	orm.RunSyncdb("default", false, true)
}
