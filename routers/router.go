package routers

import (
	"github.com/astaxie/beego"
	"github.com/rexchen123/lovelydata/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.CrawlController{})
}
