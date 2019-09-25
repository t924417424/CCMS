package controllers

import (
	"CCMS_O/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
)

func init() {
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store) //一定要写在构造函数里面，要不然第一次打开页面有可能是X
	cpt.ChallengeNums = 4
	cpt.StdWidth = 80
	cpt.StdHeight = 40
	cpt.FieldCaptchaName = "captcha"
	cpt.FieldIDName = "captchas"
}

type Base struct {
	beego.Controller
}
func (c *Base) Prepare(){
	if !models.Mysqlconn(){
		c.Ctx.Redirect(302,"/install")
	}else{
		o := orm.NewOrm()
		config := models.Config{}
		config.Id = 1
		o.Read(&config)
		c.Data["website"] = config.Website
	}
	if c.Ctx.Input.Method() == "POST" {
		getxsrf := c.GetString("_xsrf")
		xsrftoken := c.XSRFToken()
		if getxsrf != xsrftoken{
			jsonre := Newjson()
			jsonre.Code = 0
			jsonre.Msg = "页面超时，请刷新重试！"
			c.Data["json"] = &jsonre
			c.ServeJSON()
		}
	}
}
