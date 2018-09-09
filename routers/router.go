package routers

import (
	"ababila/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/test/([0-9]+):id", &controllers.TestController{})
	beego.Router("/post", &controllers.PostController{}, "post:Post")

}
