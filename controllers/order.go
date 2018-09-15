package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ababila/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type OrderController struct {
	beego.Controller
}

func (c *OrderController) respdata(resp *map[string]interface{}) {
	c.Data["json"] = resp
	c.ServeJSON()
}
func (c *OrderController) PostOrder() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	req := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	start, err := time.ParseInLocation("2006-01-02", req["start_date"].(string), time.Local)
	if err != nil {
		beego.Info(err)
	}
	end, err := time.ParseInLocation("2006-01-02", req["end_date"].(string), time.Local)
	if err != nil {
		beego.Info(err)
	}
	if bool := end.After(start); !bool {
		beego.Info("time", bool)
	}

	order := []models.OrderHouse{}
	qs := o.QueryTable("OrderHouse")
	_, err1 := qs.Filter("House__Id", req["house_id"]).RelatedSel().All(&order)
	if err1 != nil {
		beego.Info(err1)
	}
	for _, val := range order {
		if val.Status == "WAIT_ACCEPT" || val.Status == "WAIT_COMMENT" {
			bool1 := val.Begin_date.Before(start)
			bool2 := val.End_date.Before(start)
			bool3 := val.Begin_date.After(end)
			bool4 := val.End_date.After(end)
			if (bool1 && bool2) || (bool3 && bool4) {
				beego.Info("time ok")
			} else {
				resp["errno"] = models.RECODE_DATAERR
				resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
				return
			}
		}
		beego.Info("order ok")
	}

	house := models.House{}
	house.Id, _ = strconv.Atoi(req["house_id"].(string))
	err = o.Read(&house)
	if err != nil {
		beego.Info(err)
		return
	}
	_, err2 := o.LoadRelated(&house, "User")
	if err2 != nil {
		beego.Info(err)
		return
	}
	if house.User.Id == c.GetSession("user_id") {
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}

	data := models.OrderHouse{}
	uid := models.User{Id: c.GetSession("user_id").(int)}
	data.User = &uid
	hid := models.House{Id: house.Id}
	data.House = &hid
	data.Begin_date = start
	data.End_date = end
	data.Days = int(end.Sub(start).Hours() / 24)
	data.House_price = house.Price
	data.Amount = data.Days * data.House_price
	data.Status = "WAIT_ACCEPT"
	_, err3 := o.Insert(&data)
	if err3 != nil {
		beego.Info(err3)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &data.Id
}

func (c *OrderController) PutOstatus() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	req := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}

	oid := c.Ctx.Input.Param(":oid")
	order := models.OrderHouse{}
	order.Id, _ = strconv.Atoi(oid)

	if req["action"] == "accept" {
		order.Status = "WAIT_COMMENT"
		_, err := o.Update(&order, "Status")
		if err != nil {
			beego.Info(err)
			return
		}
	}
	if req["action"] == "reject" {
		order.Status = "REJECTED"
		order.Comment = req["reason"].(string)
		_, err := o.Update(&order, "Status", "Comment")
		if err != nil {
			beego.Info(err)
			return
		}
	}

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

func (c *OrderController) PutOcomment() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	req := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}

	oid := c.Ctx.Input.Param(":oid")
	order := models.OrderHouse{}
	order.Id, _ = strconv.Atoi(oid)
	order.Status = "COMPLETE"
	order.Comment = req["comment"].(string)
	_, err := o.Update(&order, "Status", "Comment")
	if err != nil {
		beego.Info(err)
		return
	}

	err1 := o.Read(&order)
	if err1 != nil {
		beego.Info(err1)
	}
	_, err2 := o.LoadRelated(&order, "House")
	if err2 != nil {
		beego.Info(err2)
	}
	house := models.House{}
	house.Id = order.House.Id
	house.Order_count = order.House.Order_count + 1
	_, err3 := o.Update(&house, "Order_count")
	if err3 != nil {
		beego.Info(err3)
		return
	}

	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}
