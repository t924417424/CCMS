package controllers

import (
	"os"
	"fmt"
	"log"
	"runtime"
	"strings"
	"CCMS_O/models"
	"encoding/json"
	"path/filepath"
	"CCMS_O/function"
	"github.com/jmoiron/sqlx"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)
type Install struct {
	beego.Controller
}
type Conn struct {
	beego.Controller
}
func (c *Install) Get(){
	step := c.GetString(":step")
	fmt.Println(step)
	if step == ""{
		c.TplName = "install/install_1.html";
	}else if step == "2" {
		c.Data["host"] = c.Ctx.Request.Host
		c.Data["path"] = getCurrentDirectory()//获取当前目录
		c.Data["var"] = "v1.0.1"
		c.Data["os"] = runtime.GOOS + " " + runtime.GOARCH
		c.Data["cpu"] = runtime.GOMAXPROCS(0)
		c.TplName = "install/install_2.html";
	}else if step == "3" {
		c.TplName = "install/install_3.html";
	}else if step == "4" {
		c.TplName = "install/install_4.html";
	}

}

func (c *Conn) Post(){
	type Configini struct {
		Dbhost string `json:"dbhost"`
		Dbuser string `json:"dbuser"`
		Dbname string `json:"dbname"`
		Dbpass string `json:"dbpass"`
	}
	var configini Configini
	//if(dbhost && dbuser && dbpas && dbname && website && username && password && salt){
	dbhost := c.GetString("dbhost")
	dbuser := c.GetString("dbuser")
	dbname := c.GetString("dbname")
	dbpass := c.GetString("dbpass")
	website := c.GetString("website")
	username := c.GetString("username")
	password := c.GetString("password")
	salt := c.GetString("salt")
	jsonre := Newjson()
	var Db *sqlx.DB
	database, _ := sqlx.Open("mysql", dbuser + ":" + dbpass + "@tcp(" + dbhost + ":3306)/")
	Db = database
	Db.Exec("CREATE DATABASE " + dbname)
	Db.Close()
	err := orm.RegisterDataBase("default","mysql",dbuser + ":" + dbpass + "@tcp(" + dbhost + ":3306)/"+ dbname + "?charset=utf8",30)
	if  err != nil{
		jsonre.Code = 0
		jsonre.Msg = "数据库连接失败！"
	}else{
		configini.Dbhost = dbhost
		configini.Dbname = dbname
		configini.Dbuser = dbuser
		configini.Dbpass = dbpass
		data, err := json.Marshal(configini)
		if err == nil{
			function.WriteWithIoutil("./conf/config.ini",string(data))
		}
		//orm.RegisterModel(new(models.Admin),new(models.Config),new(models.Dataset),new(models.BindUid))
		orm.RunSyncdb("default", true, true)
		o := orm.NewOrm()
		//o.Using("install")
		salt = function.Md5V(salt)
		password = function.Md5V(salt + password)
		admin := models.Admin{}
		admin.Username = username
		admin.Password = password
		config := models.Config{}
		config.Website = website
		config.Salt = salt
		config.Apilink = "https://wauth.liuxiaoa.com/qqlogin.php"
		config.Apikey = "123456"
		dataset := models.Dataset{}
		o.Insert(&admin)
		o.Insert(&config)
		o.Insert(&dataset)
		jsonre.Code = 1
		jsonre.Msg = "数据库连接成功！开始安装数据。。。"
	}
	c.Data["json"] = jsonre
	c.ServeJSON()
}


func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}