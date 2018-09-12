package routers

import (
	"ababila/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetArea")
	beego.Router("/api/v1.0/houses/index", &controllers.HouseController{}, "get:GetHid")
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:GetSes;delete:DelSes")
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:PostUser")
	beego.Router("/api/v1.0/sessions", &controllers.SessionController{}, "post:PostSes")
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{}, "post:PostAvatar")
	beego.Router("/api/v1.0/user", &controllers.UserController{}, "get:GetUser")
	beego.Router("/api/v1.0/user/name", &controllers.UserController{}, "put:PutUname")
	beego.Router("/api/v1.0/user/auth", &controllers.UserController{}, "get:GetUser;post:PostUid")
	beego.Router("/api/v1.0/user/houses", &controllers.UserController{}, "get:GetHouse")
	beego.Router("/api/v1.0/houses", &controllers.HouseController{}, "post:PostHouse")

}
