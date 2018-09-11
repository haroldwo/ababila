package controllers

import (
	"ababila/models"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
)

type AreaController struct {
	beego.Controller
}

func (c *AreaController) respdata(resp *map[string]interface{}) {
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *AreaController) GetArea() {
	beego.Info("connect success")
	o := orm.NewOrm()
	var area []*models.Area
	resp := make(map[string]interface{})
	defer c.respdata(&resp)

	client := redis.NewClient(&redis.Options{
		Addr:     "172.16.0.10:6379",
		Password: "fred",
		DB:       0,
	})
	value, err := client.Get("area").Bytes()
	if err != nil {
		beego.Info(err)
	}
	if value != nil {
		err := json.Unmarshal(value, &area)
		if err != nil {
			beego.Info(err)
			resp["errno"] = models.RECODE_DATAERR
			resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
			return
		}
		resp["errno"] = models.RECODE_OK
		resp["errmsg"] = models.RecodeText(models.RECODE_OK)
		resp["data"] = &area
		beego.Info("use cache")
		return
	}

	qs := o.QueryTable("area")
	num, err := qs.All(&area)
	if err != nil {
		beego.Info(err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	if num == 0 {
		beego.Info("no data")
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &area

	areajson, err := json.Marshal(&area)
	if err != nil {
		beego.Info(err)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	err1 := client.Set("area", areajson, time.Second*3600)
	if err1 != nil {
		beego.Info(err1)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
}
