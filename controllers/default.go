package controllers

import (
	"CCMS_O/models"
	"github.com/astaxie/beego/orm"
)


type MainController struct {
	Base
}

func (c *MainController) Get() {
	set := models.Dataset{}
	set.Id = 1
	o := orm.NewOrm()
	o.Read(&set)
	c.Data["onlytitle"] = set.Onlytitle
	c.TplName = "Index/index.html"
}
