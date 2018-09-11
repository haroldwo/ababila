package main

import (
	_ "ababila/models"
	_ "ababila/routers"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

//func insertData() {
//	o := orm.NewOrm()
//	o.Using("default")
//
//	user := models.User{}
//	user.Name = "fred"
//	id, err := o.Insert(&user)
//	if err != nil {
//		beego.Info("insert err")
//		return
//	}
//	beego.Info("insert success", id)
//}
//
//func queryData() {
//	o := orm.NewOrm()
//	o.Using("default")
//
//	user := models.User{Id: 1}
//	err := o.Read(&user)
//	if err != nil {
//		beego.Info("query err")
//		return
//	}
//	beego.Info("query success", user)
//}
//
//func updateData() {
//	o := orm.NewOrm()
//	o.Using("default")
//
//	user := models.User{Id: 1, Name: "harold"}
//	_, err := o.Update(&user)
//	if err != nil {
//		beego.Info("update err")
//		return
//	}
//	beego.Info("update success", user)
//}
//
//func deleteData() {
//	o := orm.NewOrm()
//	o.Using("default")
//
//	user := models.User{Id: 1}
//	_, err := o.Delete(&user)
//	if err != nil {
//		beego.Info("delete err")
//		return
//	}
//	beego.Info("delete success", user)
//}
//func insertOrder() {
//	o := orm.NewOrm()
//	o.Using("default")
//
//	data := models.User_order{}
//	data.Order_data = "burgerking"
//	user := models.User{Id: 1}
//	data.User = &user
//	id, err := o.Insert(&data)
//	if err != nil {
//		beego.Info("insert err")
//		return
//	}
//	beego.Info("insert success", id)
//}
//
//func queryOrder() {
//	var orders []*models.User_order
//	o := orm.NewOrm()
//	qs := o.QueryTable("User_order")
//	_, err := qs.Filter("User__Id", 1).All(&orders)
//	if err != nil {
//		beego.Info("query err")
//		return
//	}
//	for _, v := range orders {
//		beego.Info("query order", v)
//	}
//}

func ignoreStaticPath() {

	//透明static
	//beego.SetStaticPath("group1/M00/", "fdfs/storage_data/data/")

	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
}

func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url: ", orpath)
	//如果请求uri还有api字段,说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
}

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	ignoreStaticPath()
	beego.Run()

}
