package extensions

import (
	"github.com/astaxie/beego"
)

func Init() {
	conf := beego.AppConfig
	initMongo(conf)
	initQueue()
}
