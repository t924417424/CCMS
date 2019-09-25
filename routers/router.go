package routers

import (
	"strings"
	"CCMS_O/function"
	"CCMS_O/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	var FilterUser = func(ctx *context.Context) {
		_, ok := ctx.Input.Session("suid").(int)
		ok2 := strings.Contains(ctx.Request.RequestURI, "/ccms/c_login")
		if !ok && !ok2 {
			ctx.Redirect(302, "/")
		}
	}

	beego.InsertFilter("/ccms/*",beego.BeforeRouter,FilterUser)

	var FilterLock = func(ctx *context.Context) {
		ok := function.IsExist("./conf/config.ini")
		ok2 := strings.Contains(ctx.Request.RequestURI, "install")
		if ok && ok2 {
			ctx.Redirect(302, "/login")
		}
	}
	var FilterInstall = func(ctx *context.Context) {
		ok := function.IsExist("./conf/config.ini")
		ok2 := strings.Contains(ctx.Request.RequestURI, "install")
		if !ok && !ok2 {
			ctx.Redirect(302, "/install")
		}
	}
	beego.InsertFilter("/install",beego.BeforeRouter,FilterLock)
	beego.InsertFilter("/*",beego.BeforeRouter,FilterInstall)

    beego.Router("/", &controllers.MainController{})
    beego.Router("/install/?:step",&controllers.Install{})
	beego.Router("/install/conn",&controllers.Conn{})
    beego.Router("/login",&controllers.Login{})
    beego.Router("/qqlogin",&controllers.QQLogin{})
    beego.Router("/ccms/c_login",&controllers.Clogin{})
    beego.Router("/ccms/index",&controllers.Cindex{})
	beego.Router("/ccms/left.html",&controllers.Left{})
	beego.Router("/ccms/top.html",&controllers.Top{})
    beego.Router("/ccms/main",&controllers.Main{})
    beego.Router("/ccms/upload",&controllers.Cupload{})
	beego.Router("/ccms/getexcel", &controllers.Cgetexcel{})
    beego.Router("/ccms/sysset",&controllers.Sysset{})
    beego.Router("/ccms/importexcel",&controllers.Importdata{})
    beego.Router("/ccms/datasys",&controllers.Datasys{})
    beego.Router("/ccms/repass",&controllers.Repass{})
    beego.Router("/ccms/select",&controllers.Selects{})
    beego.Router("/ccms/api/?:select",&controllers.Selects{},"*:Select")
	beego.Router("/ccms/bind",&controllers.Bind{})
	beego.Router("/ccms/bind/?:select",&controllers.Bind{},"*:Bind")
}
