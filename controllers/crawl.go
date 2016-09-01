package controllers

import (
	"github.com/astaxie/beego"
	"github.com/rexchen123/lovelydata/extensions"
)

type CrawlController struct {
	beego.Controller
}

func (c *CrawlController) Start() {
	extensions.Q.Queue.Put("eoobird")
	c.TplName = "index.tpl"
}
