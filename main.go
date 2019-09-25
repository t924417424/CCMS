package main

import (
	_ "CCMS_O/models"
	_ "CCMS_O/routers"
	"github.com/astaxie/beego"
)


func main(){
	beego.BConfig.RecoverPanic = true
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}