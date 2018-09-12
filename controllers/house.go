package controllers

import (
	"ababila/models"
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type HouseController struct {
	beego.Controller
}

func (c *HouseController) respdata(resp *map[string]interface{}) {
	c.Data["json"] = resp
	c.ServeJSON()
}
func (c *HouseController) GetHid() {
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	beego.Info("data err")
	resp["errno"] = models.RECODE_DBERR
	resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
}

func (c *HouseController) PostHouse() {
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
	if m, _ := regexp.MatchString("^.{0,20}$", resp["title"].(string)); !m {
		beego.Info("regexp", resp["title"].(string))
		return
	}
	if m, _ := regexp.MatchString("^\\d{1,10}$", resp["price"].(string)); !m {
		beego.Info("regexp", resp["price"].(int))
		return
	}
	if m, _ := regexp.MatchString("^.{0,100}$", resp["address"].(string)); !m {
		beego.Info("regexp", resp["address"].(string))
		return
	}
	if m, _ := regexp.MatchString("^[1-9]$", resp["room_count"].(string)); !m {
		beego.Info("regexp", resp["room_count"].(int))
		return
	}
	if m, _ := regexp.MatchString("^\\d{1,3}$", resp["acreage"].(string)); !m {
		beego.Info("regexp", resp["acreage"].(int))
		return
	}
	if m, _ := regexp.MatchString("^.{0,20}$", resp["unit"].(string)); !m {
		beego.Info("regexp", resp["unit"].(string))
		return
	}
	if m, _ := regexp.MatchString("^[1-9]$", resp["capacity"].(string)); !m {
		beego.Info("regexp", resp["capacity"].(int))
		return
	}
	if m, _ := regexp.MatchString("^.{0,20}$", resp["beds"].(string)); !m {
		beego.Info("regexp", resp["beds"].(string))
		return
	}
	if m, _ := regexp.MatchString("^\\d{0,5}$", resp["deposit"].(string)); !m {
		beego.Info("regexp", resp["deposit"].(int))
		return
	}
	if m, _ := regexp.MatchString("^[1-9]\\d{0,1}$", resp["min_days"].(string)); !m {
		beego.Info("regexp", resp["min_days"].(int))
		return
	}
	if m, _ := regexp.MatchString("^[1-9]\\d{0,1}$", resp["max_days"].(string)); !m {
		beego.Info("regexp", resp["max_days"].(int))
		return
	}

	data := models.House{}
	user := models.User{Id: c.GetSession("user_id").(int)}
	err1 := o.Read(&user)
	if err1 != nil {
		beego.Info(err1)
	}
	data.User = &user
	str, _ := strconv.Atoi(resp["area_id"].(string))
	area := models.Area{Id: str}
	err2 := o.Read(&area)
	if err2 != nil {
		beego.Info(err2)
	}
	data.Area = &area
	data.Title = resp["title"].(string)
	data.Price, _ = strconv.Atoi(resp["price"].(string))
	data.Address = resp["address"].(string)
	data.Room_count, _ = strconv.Atoi(resp["room_count"].(string))
	data.Acreage, _ = strconv.Atoi(resp["acreage"].(string))
	data.Unit = resp["unit"].(string)
	data.Capacity, _ = strconv.Atoi(resp["capacity"].(string))
	data.Beds = resp["beds"].(string)
	data.Deposit, _ = strconv.Atoi(resp["deposit"].(string))
	data.Min_days, _ = strconv.Atoi(resp["min_days"].(string))
	data.Max_days, _ = strconv.Atoi(resp["max_days"].(string))
	hid, err := o.Insert(&data)
	if err != nil {
		beego.Info(err)
		return
	}

	house := models.House{Id: int(hid)}
	m2m := o.QueryM2M(&house, "Facilities")
	for _, v := range resp["facility"].([]interface{}) {
		id, _ := strconv.Atoi(v.(string))
		facilities := &models.Facility{Id: id}
		o.Insert(facilities)
		_, err := m2m.Add(facilities)
		if err != nil {
			beego.Info(err)
		}
	}

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = hid
}
