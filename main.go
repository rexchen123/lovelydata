package main

import (
	"github.com/astaxie/beego"
	"github.com/rexchen123/lovelydata/extensions"
	"github.com/rexchen123/lovelydata/models"
	_ "github.com/rexchen123/lovelydata/routers"
)

func main() {
	beego.SetLogger("console", "")
	extensions.Init()
	models.WaitToCrawlUser()
	beego.Run()
}
