package controllers

import (
	"ababila/models"
	"encoding/json"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
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
	req := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	if m, _ := regexp.MatchString("^.{0,20}$", req["title"].(string)); !m {
		beego.Info("regexp", req["title"].(string))
		return
	}
	if m, _ := regexp.MatchString("^\\d{1,10}$", req["price"].(string)); !m {
		beego.Info("regexp", req["price"].(string))
		return
	}
	if m, _ := regexp.MatchString("^.{0,100}$", req["address"].(string)); !m {
		beego.Info("regexp", req["address"].(string))
		return
	}
	if m, _ := regexp.MatchString("^[1-9]$", req["room_count"].(string)); !m {
		beego.Info("regexp", req["room_count"].(string))
		return
	}
	if m, _ := regexp.MatchString("^\\d{1,3}$", req["acreage"].(string)); !m {
		beego.Info("regexp", req["acreage"].(string))
		return
	}
	if m, _ := regexp.MatchString("^.{0,20}$", req["unit"].(string)); !m {
		beego.Info("regexp", req["unit"].(string))
		return
	}
	if m, _ := regexp.MatchString("\\d{1,2}$", req["capacity"].(string)); !m {
		beego.Info("regexp", req["capacity"].(string))
		return
	}
	if m, _ := regexp.MatchString("^.{0,20}$", req["beds"].(string)); !m {
		beego.Info("regexp", req["beds"].(string))
		return
	}
	if m, _ := regexp.MatchString("^\\d{0,5}$", req["deposit"].(string)); !m {
		beego.Info("regexp", req["deposit"].(string))
		return
	}
	if m, _ := regexp.MatchString("^[1-9]\\d{0,1}$", req["min_days"].(string)); !m {
		beego.Info("regexp", req["min_days"].(string))
		return
	}
	if m, _ := regexp.MatchString("^\\d{1,2}$", req["max_days"].(string)); !m {
		beego.Info("regexp", req["max_days"].(string))
		return
	}

	data := models.House{}
	user := models.User{Id: c.GetSession("user_id").(int)}
	data.User = &user
	str, _ := strconv.Atoi(req["area_id"].(string))
	area := models.Area{Id: str}
	data.Area = &area
	data.Title = req["title"].(string)
	data.Price, _ = strconv.Atoi(req["price"].(string))
	data.Address = req["address"].(string)
	data.Room_count, _ = strconv.Atoi(req["room_count"].(string))
	data.Acreage, _ = strconv.Atoi(req["acreage"].(string))
	data.Unit = req["unit"].(string)
	data.Capacity, _ = strconv.Atoi(req["capacity"].(string))
	data.Beds = req["beds"].(string)
	data.Deposit, _ = strconv.Atoi(req["deposit"].(string))
	data.Min_days, _ = strconv.Atoi(req["min_days"].(string))
	data.Max_days, _ = strconv.Atoi(req["max_days"].(string))
	facilities := []models.Facility{}
	for _, v := range req["facility"].([]interface{}) {
		id, _ := strconv.Atoi(v.(string))
		fmap := models.Facility{Id: id}
		facilities = append(facilities, fmap)
	}

	hid, err := o.Insert(&data)
	if err != nil {
		beego.Info(err)
		return
	}
	house := models.House{Id: int(hid)}
	m2m := o.QueryM2M(&house, "Facilities")
	num, err := m2m.Add(facilities)
	if err != nil || num == 0 {
		beego.Info("m2m.Add", err)
	}
	hmap := make(map[string]interface{})
	hmap["house_id"] = hid
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = hmap
}

func (c *HouseController) PostHimage() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}

	file, fhead, err1 := c.GetFile("house_image")
	if err1 != nil {
		beego.Info(err1)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	defer file.Close()
	if m, err := regexp.MatchString("^\\.(png|jpg|bmp)$", path.Ext(fhead.Filename)); !m {
		beego.Info("regexp", err)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	ffield := "static/upload/" + fhead.Filename
	err2 := c.SaveToFile("house_image", ffield)
	if err2 != nil {
		beego.Info(err2)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}

	id := c.Input().Get("house_id")
	hid, err5 := strconv.Atoi(id)
	if err5 != nil {
		beego.Info(err5)
		return
	}
	house := models.House{}
	house.Id = hid
	err3 := o.Read(&house)
	if err3 != nil {
		beego.Info(err3)
		return
	}
	if house.Index_image_url == "" {
		house.Index_image_url = ffield
		_, err := o.Update(&house, "index_image_url")
		if err != nil {
			beego.Info(err)
			return
		}
	}
	himage := models.HouseImage{}
	himage.House = &house
	himage.Url = ffield
	_, err4 := o.Insert(&himage)
	if err4 != nil {
		beego.Info(err4)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	data := make(map[string]string)
	data["url"] = "http://localhost:8080/" + himage.Url
	resp["data"] = &data
}

func (c *HouseController) GetHouse() {
	resp := make(map[string]interface{})
	defer c.respdata(&resp)
	o := orm.NewOrm()
	client := redis.NewClient(&redis.Options{
		Addr:     "172.16.0.10:6379",
		Password: "fred",
		DB:       0,
	})
	if c.GetSession("user_id") == nil {
		beego.Info("Session", c.GetSession("user_id"))
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}

	hid := c.Ctx.Input.Param(":hid")
	value, err := client.Get(hid).Bytes()
	if err != nil {
		beego.Info(err)
	} else if value != nil {
		hmap := make(map[string]interface{})
		err := json.Unmarshal(value, &hmap)
		if err != nil {
			beego.Info(err)
			return
		}
		resp["errno"] = models.RECODE_OK
		resp["errmsg"] = models.RecodeText(models.RECODE_OK)
		resp["data"] = &hmap
		beego.Info("use redis")
		return
	}

	house := models.House{}
	house.Id, err = strconv.Atoi(hid)
	if err != nil {
		beego.Info(err)
		return
	}
	err = o.Read(&house)
	if err != nil {
		beego.Info(err)
		return
	}
	//_, err1 := o.LoadRelated(&house, "OrderHouse")
	//if err1 != nil {
	//	beego.Info(err)
	//	return
	//}
	_, err1 := o.LoadRelated(&house, "User")
	if err1 != nil {
		beego.Info(err)
		return
	}
	_, err1 = o.LoadRelated(&house, "Facilities")
	if err1 != nil {
		beego.Info(err)
		return
	}
	_, err1 = o.LoadRelated(&house, "Images")
	if err1 != nil {
		beego.Info(err)
		return
	}

	data := make(map[string]interface{})
	data["acreage"] = house.Acreage
	data["beds"] = house.Beds
	data["capacity"] = house.Capacity
	data["address"] = house.Address
	data["deposit"] = house.Deposit
	data["min_days"] = house.Min_days
	data["max_days"] = house.Max_days
	data["unit"] = house.Unit
	data["price"] = house.Price
	data["room_count"] = house.Room_count
	data["title"] = house.Title
	data["user_id"] = house.User.Id
	data["user_name"] = house.User.Name
	data["user_avatar"] = house.User.Avatar_url
	hfacility := []interface{}{}
	for _, val := range house.Facilities {
		hfacility = append(hfacility, val.Id)
	}
	data["facilities"] = &hfacility
	himages := []interface{}{}
	for _, val := range house.Images {
		himages = append(himages, val.Url)
	}
	data["img_urls"] = &himages
	horder := []interface{}{}
	for _, val := range house.Orders {
		odata := make(map[string]interface{})
		odata["user_name"] = val.User.Name
		odata["comment"] = val.Comment
		odata["ctime"] = val.Ctime
		horder = append(horder, odata)
	}
	data["comments"] = &horder

	hmap := make(map[string]interface{})
	hmap["house"] = &data
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &hmap

	hjson, err := json.Marshal(&hmap)
	if err != nil {
		beego.Info(err)
		return
	}
	rds := client.Set(hid, hjson, time.Second*60)
	beego.Info(rds)
}
