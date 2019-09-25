package models

import (
	"encoding/json"
	"CCMS_O/function"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)
type Admin struct {
	Id   int `orm:"unique"`
	Username string `orm:"size(32)"`
	Password string	`orm;"size(32)"`
	Logintime string `orm:"size(13)"`
}
type Config struct {
	Id int `orm:"unique"`
	Website string	`orm:"size(200)"`
	Salt	string	`orm:"size(32);null"`
	Mode int `orm:"default(1);description(当前选择的登陆模式)"`
	Apilink string `orm:"default(https://wauth.liuxiaoa.com/qqlogin.php)"`
	//用于获取QQ ACC_TOKEN的接口地址
	Apikey string `orm:"default(123456)"`
	//用于获取QQ ACC_TOKEN的接口地址的操作密钥
}
type Dataset struct {
	Id int `orm:"unique"`
	Exceltitle string
	Looktitle string	//用户登陆之后可以展示的字段
	Onlytitle string `orm:"size(50);description(表格字段中的唯一不重复字段，用于登陆绑定)"`
	Passtitle string `orm:"size(50);description(表格字段中的密码字段，用于登陆验证)"`
}
type BindUid struct {
	Id int `orm:"unique"`
	//Onlytitle string `orm:"size(50)"`
	Binddata string `orm:"description(绑定的数据，可以是QQUID也可以是用户密码)"`
	Binduser string `orm:"size(50);unique"`
	Line int `orm:description(在Excel中对应的行数，避免Redis中未缓存导致的重复查询，如与Excel中的唯一字段不匹配则更新！)`
	Logintime string `orm:"size(13)"`
}
var sqlstatus = 0
func init(){
	orm.Debug = false
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//orm.RegisterDataBase("default","mysql","root:root@/testgo?charset=utf8",30)

	orm.RegisterModel(new(Admin),new(Config),new(Dataset),new(BindUid))
	if Mysqlconn(){
		orm.RunSyncdb("default", false, true)
	}
}
func Mysqlconn() bool{
	type Configini struct {
		Dbhost string `json:"dbhost"`
		Dbuser string `json:"dbuser"`
		Dbname string `json:"dbname"`
		Dbpass string `json:"dbpass"`
	}
	if sqlstatus == 0{
		if function.IsExist("./conf/config.ini"){
			var config  Configini
			data,err := function.ReadAllIntoMemory("./conf/config.ini")
			if err != nil{
				return false
			}else{
				json.Unmarshal([]byte(data), &config)
				orm.RegisterDriver("mysql", orm.DRMySQL)
				orm.RegisterDataBase("default","mysql",config.Dbuser + ":" + config.Dbpass + "@tcp(" + config.Dbhost + ":3306)/"+ config.Dbname + "?charset=utf8",30)
				sqlstatus = 1
				return true
			}
		}else{
			return false
		}
	}else{
		return true
	}
}

