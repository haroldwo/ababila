package controllers

import (
	"ababila/models"
	"encoding/json"
	"path"
	"regexp"
	"sort"
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

func (c *HouseController) GetHindex() {
	o := orm.NewOrm()
	resp := make(map[string]interface{})
	defer c.respdata(&resp)

	house := []models.House{}
	qs := o.QueryTable("House")
	_, err := qs.All(&house)
	if err != nil {
		beego.Info(err)
		return
	}
	hindex := []interface{}{}
	for _, val := range house {
		data := make(map[string]interface{})
		data["house_id"] = val.Id
		data["img_url"] = val.Index_image_url
		data["title"] = val.Title
		hindex = append(hindex, data)
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &hindex
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
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^\\d{1,10}$", req["price"].(string)); !m {
		beego.Info("regexp", req["price"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^.{0,100}$", req["address"].(string)); !m {
		beego.Info("regexp", req["address"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^[1-9]$", req["room_count"].(string)); !m {
		beego.Info("regexp", req["room_count"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^\\d{1,3}$", req["acreage"].(string)); !m {
		beego.Info("regexp", req["acreage"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^.{0,20}$", req["unit"].(string)); !m {
		beego.Info("regexp", req["unit"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("\\d{1,2}$", req["capacity"].(string)); !m {
		beego.Info("regexp", req["capacity"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^.{0,20}$", req["beds"].(string)); !m {
		beego.Info("regexp", req["beds"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^\\d{0,5}$", req["deposit"].(string)); !m {
		beego.Info("regexp", req["deposit"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^[1-9]\\d{0,1}$", req["min_days"].(string)); !m {
		beego.Info("regexp", req["min_days"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if m, _ := regexp.MatchString("^\\d{1,2}$", req["max_days"].(string)); !m {
		beego.Info("regexp", req["max_days"].(string))
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
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

	test := models.OrderHouse{}
	uid := models.User{Id: c.GetSession("user_id").(int)}
	test.User = &uid
	id := models.House{Id: int(hid)}
	test.House = &id
	test.Begin_date = time.Now()
	test.End_date = time.Now()
	test.Days = 0
	test.House_price, _ = strconv.Atoi(req["price"].(string))
	test.Amount = test.Days * test.House_price
	test.Status = "TEST"
	_, err1 := o.Insert(&test)
	if err1 != nil {
		beego.Info(err1)
		resp["errno"] = models.RECODE_DATAERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
	}
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
	_, err4 := o.LoadRelated(&house, "Orders")
	if err4 != nil {
		beego.Info(err4)
	}
	_, err1 := o.LoadRelated(&house, "User")
	if err1 != nil {
		beego.Info(err1)
		return
	}
	_, err2 := o.LoadRelated(&house, "Facilities")
	if err2 != nil {
		beego.Info(err2)
		return
	}
	_, err3 := o.LoadRelated(&house, "Images")
	if err3 != nil {
		beego.Info(err3)
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
	data["hid"] = house.Id
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

	qs := o.QueryTable("OrderHouse")
	_, err5 := qs.Filter("House__Id", house.Id).RelatedSel().All(&house.Orders)
	if err5 != nil {
		beego.Info(err1)
	}
	horder := []interface{}{}
	for _, val := range house.Orders {
		if val.Status != "COMPLETE" {
			continue
		}
		odata := make(map[string]interface{})
		odata["user_name"] = val.User.Name
		odata["comment"] = val.Comment
		odata["ctime"] = val.Ctime.Format("2006-01-01 15:04:05")
		horder = append(horder, odata)
	}
	data["comments"] = &horder

	hmap := make(map[string]interface{})
	hmap["house"] = &data
	hmap["user_id"] = c.GetSession("user_id")
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

func (c *HouseController) GetHsearch() {
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

	aid, _ := strconv.Atoi(c.Input().Get("aid"))
	start, err := time.ParseInLocation("2006-01-02", c.Input().Get("sd"), time.Local)
	if err != nil {
		beego.Info(err)
	}
	end, err := time.ParseInLocation("2006-01-02", c.Input().Get("ed"), time.Local)
	if err != nil {
		beego.Info(err)
	}
	hid := c.Checkdata(start, end)
	house := models.House{}
	hdata := make(map[string]interface{})
	data := []interface{}{}
	for _, val := range hid {
		house.Id = val
		o.Read(&house)
		beego.Info("here1")
		if house.Area.Id == aid {
			beego.Info("here2")
			_, err1 := o.LoadRelated(&house, "User")
			beego.Info("here3")
			if err1 != nil {
				beego.Info(err1)
				return
			}
			hdata["house_id"] = house.Id
			hdata["img_url"] = house.Index_image_url
			hdata["user_avatar"] = house.User.Avatar_url
			hdata["price"] = house.Price
			hdata["title"] = house.Title
			hdata["room_count"] = house.Room_count
			hdata["order_count"] = house.Order_count
			hdata["address"] = house.Address
			data = append(data, hdata)
			beego.Info("here4")
		}
	}
	if &data == nil {
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	smap := make(map[string]interface{})
	smap["houses"] = &data
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &smap
}

func (c *HouseController) Checkdata(start time.Time, end time.Time) []int {
	o := orm.NewOrm()
	order := []models.OrderHouse{}
	hid := []int{}
	qs := o.QueryTable("OrderHouse")
	_, err1 := qs.RelatedSel().All(&order)
	if err1 != nil {
		beego.Info(err1)
	}
	for _, val := range order {
		if val.Status == "" {
			return nil
		}
		bool1 := val.Begin_date.Before(start)
		bool2 := val.End_date.Before(start)
		bool3 := val.Begin_date.After(end)
		bool4 := val.End_date.After(end)
		if (bool1 && bool2) || (bool3 && bool4) {
			beego.Info("time ok")
			hid = append(hid, val.House.Id)
		}
	}
	sort.Ints(hid)
	return RemoveDuplicatesAndEmpty(hid)
}

func RemoveDuplicatesAndEmpty(a []int) (ret []int) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if i > 0 && a[i-1] == a[i] {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
