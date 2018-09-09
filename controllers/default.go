package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}
type TestController struct {
	beego.Controller
}
type PostController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *TestController) Get() {
	id := c.Ctx.Input.Param(":id")
	c.Ctx.WriteString("get id:" + id)
}
func (c *PostController) Post() {
	c.Ctx.WriteString("get post.")
}
