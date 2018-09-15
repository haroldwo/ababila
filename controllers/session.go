package controllers

import (
	"ababila/models"
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type SessionController struct {
	beego.Controller
}

func (c *SessionController) respdata(resp *map[string]interface{}) {
	c.Data["json"] = resp
	c.ServeJSON()
}
func (c *SessionController) GetSes() {
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	user := models.User{}

	name := c.GetSession("name")
	if name != nil {
		user.Name = name.(string)
		resp["errno"] = models.RECODE_OK
		resp["errmsg"] = models.RecodeText(models.RECODE_OK)
		resp["data"] = &user
		return
	}
	resp["errno"] = models.RECODE_SESSIONERR
	resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
}

func (c *SessionController) DelSes() {
	resp := make(map[string]interface{})
	defer c.respdata(&resp)

	c.DelSession("name")
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

func (c *SessionController) PostSes() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	json.Unmarshal(c.Ctx.Input.RequestBody, &resp)

	if resp["mobile"] == nil || resp["password"] == nil {
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	user := models.User{}
	user.Mobile = resp["mobile"].(string)
	user.Password_hash = resp["password"].(string)

	qs := o.QueryTable("user")
	err1 := qs.Filter("Mobile", user.Mobile).One(&user)
	if err1 != nil {
		resp["errno"] = models.RECODE_USERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_USERERR)
		return
	}
	if resp["password"] != user.Password_hash {
		resp["errno"] = models.RECODE_PWDERR
		resp["errmsg"] = models.RecodeText(models.RECODE_PWDERR)
		return
	}

	c.SetSession("name", user.Name)
	c.SetSession("mobile", user.Mobile)
	c.SetSession("user_id", user.Id)
	name := c.GetSession("name")
	if name != nil {
		user.Name = name.(string)
		resp["errno"] = models.RECODE_OK
		resp["errmsg"] = models.RecodeText(models.RECODE_OK)
		resp["data"] = &user
		return
	}
	resp["errno"] = models.RECODE_DBERR
	resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
}
