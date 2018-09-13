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

func (c *UserController) GetUser() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)

	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	user := models.User{}
	user.Id = c.GetSession("user_id").(int)
	err := o.Read(&user)
	if err != nil {
		beego.Info(err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &user
}

func (c *UserController) PutUname() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	json.Unmarshal(c.Ctx.Input.RequestBody, &resp)

	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	data := models.User{}
	data.Name = resp["name"].(string)
	data.Id = c.GetSession("user_id").(int)
	if m, _ := regexp.MatchString("^.{0,20}$", data.Name); !m {
		beego.Info("regexp", data.Name)
		return
	}
	qs := o.QueryTable("user")
	err1 := qs.Filter("Name", data.Name).One(&data)
	if err1 == nil {
		beego.Info(err1)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	_, err2 := o.Update(&data, "name")
	if err2 != nil {
		beego.Info(err2)
		return
	}
	c.SetSession("name", data.Name)
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &data.Name
}

func (c *UserController) PostUid() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	json.Unmarshal(c.Ctx.Input.RequestBody, &resp)

	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	data := models.User{}
	data.Real_name = resp["real_name"].(string)
	data.Id_card = resp["id_card"].(string)
	data.Id = c.GetSession("user_id").(int)
	if m, _ := regexp.MatchString("^.{0,20}$", data.Real_name); !m {
		beego.Info("regexp", data.Real_name)
		return
	}
	if m, _ := regexp.MatchString("^[1-9]\\d{5}(18|19|([23]\\d))\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$", data.Id_card); !m {
		beego.Info("regexp", data.Id_card)
		return
	}
	_, err2 := o.Update(&data, "real_name", "id_card")
	if err2 != nil {
		beego.Info(err2)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

func (c *UserController) GetHouse() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)

	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	mdhouse := []models.House{}
	uid := c.GetSession("user_id").(int)
	qs := o.QueryTable("house")
	num, err1 := qs.Filter("User__Id", uid).RelatedSel().All(&mdhouse)
	if err1 != nil || num == 0 {
		beego.Info(err1)
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	beego.Info("Query", num)
	house := []interface{}{}
	for _, val := range mdhouse {
		data := make(map[string]interface{})
		data["address"] = val.Address
		data["area_name"] = val.Area.Name
		data["ctime"] = string(val.Ctime.Format("2006-01-01 15:04:05"))
		data["house_id"] = val.Id
		data["img_url"] = val.Index_image_url
		data["order_count"] = val.Order_count
		data["price"] = val.Price
		data["room_count"] = val.Room_count
		data["title"] = val.Title
		data["user_avatar"] = val.User.Avatar_url
		house = append(house, data)
	}
	hmap := make(map[string]interface{})
	hmap["houses"] = &house
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &hmap
}
