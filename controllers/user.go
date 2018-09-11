package controllers

import (
	"ababila/models"
	"encoding/json"
	"path"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) respdata(resp *map[string]interface{}) {
	c.Data["json"] = resp
	c.ServeJSON()
}
func (c *UserController) PostUser() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	json.Unmarshal(c.Ctx.Input.RequestBody, &resp)

	data := models.User{}
	data.Name = resp["mobile"].(string)
	data.Password_hash = resp["password"].(string)
	data.Mobile = resp["mobile"].(string)

	if m, _ := regexp.MatchString("^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\\d{8}$", data.Name); !m {
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}

	if m, _ := regexp.MatchString("^[0-9A-Za-z]{4,}$", data.Password_hash); !m {
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}

	qs := o.QueryTable("user")
	err1 := qs.Filter("Name", data.Name).One(&data)
	if err1 == nil {
		resp["errno"] = models.RECODE_DATAEXIST
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAEXIST)
		return
	}

	id, err2 := o.Insert(&data)
	if err2 != nil {
		beego.Info(err2)
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	beego.Info("register success, id =", id)
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	c.SetSession("name", data.Name)
}

func (c *UserController) PostAvatar() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	f, h, err1 := c.GetFile("avatar")
	if err1 != nil {
		beego.Info(err1)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	defer f.Close()
	if m, err := regexp.MatchString("^\\.(png|jpg|bmp)$", path.Ext(h.Filename)); !m {
		beego.Info("regexp", err)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	err2 := c.SaveToFile("avatar", "static/upload/"+h.Filename)
	if err2 != nil {
		beego.Info(err2)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	user := models.User{}
	user.Id = c.GetSession("user_id").(int)
	user.Avatar_url = "static/upload/" + h.Filename
	_, err := o.Update(&user, "avatar_url")
	if err != nil {
		beego.Info(err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	data := make(map[string]string)
	data["avatar_url"] = "http://localhost:8080/" + user.Avatar_url
	resp["data"] = &data
}
