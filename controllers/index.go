package controllers

import (
	"ababila/models"

	"github.com/astaxie/beego"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) respdata(resp *map[string]interface{}) {
	c.Data["json"] = resp
	c.ServeJSON()
}
func (c *IndexController) GetIndex() {
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	beego.Info("data err")
	resp["errno"] = models.RECODE_DBERR
	resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
}
