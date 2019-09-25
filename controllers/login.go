package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"CCMS_O/models"
	"CCMS_O/function"
	"github.com/tealeg/xlsx"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Login struct {
	Base
}

type QQLogin struct {
	beego.Controller
}

type Jsonre struct {
	Code int `json:"code"`
	Msg  string	`json:"msg"`
	Data string `json:"data"`
}
func Newjson() *Jsonre {
	return &Jsonre{}
}
/*
func (c *Login) Get(){
	set := models.Dataset{}
	set.Id = 1
	o := orm.NewOrm()
	o.Read(&set)
	c.Data["onlytitle"] = set.Onlytitle
	c.TplName = "Index/index.html"
}
*/
func (c *Login) Post(){
	var pasjc string
	username := c.GetString("username")
	password := c.GetString("password")
	jsonre := Newjson()
	jsonre.Code = 0
	if !cpt.VerifyReq(c.Ctx.Request) {
		jsonre.Msg = "验证码错误"
	} else if username == ""||password == ""{
		jsonre.Msg = "参数不能为空"
	}else{
		o := orm.NewOrm()
		set := models.Dataset{}
		set.Id = 1
		o.Read(&set)
		onlytitle := set.Onlytitle
		passtitle := set.Passtitle
		looktitle := set.Looktitle
		if onlytitle == "" || passtitle == ""{
			jsonre.Msg = "数据未初始化，请联系网站管理员！"
		}else{
			userbind := models.BindUid{}
			userbind.Binduser = username
			o.Read(&userbind,"binduser")	//使用指定列查找

			path := "./filetmp/" + "c48a93c545c647dc9c20a4420e6eee33"
			text, _ := function.ReadAllIntoMemory(path)
			key := []byte("scoresdcet111246141score")
			x2 := function.Decrypt3DES(text, key)
			xlFile, err := xlsx.OpenBinary(x2)
			sheet := xlFile.Sheets[0]	//选定Sheet工作表
			celid := 0	//唯一索引在EXCELTITLE中对应的列数
			pasid := 0  //密码索引在EXCELTITLE中对应的列数
			if err != nil {
				jsonre.Msg = "数据源读取出错，请联系管理员！"
			}else{
				titledata := sheet.Rows[0].Cells
				fmt.Println("标题行数据：",titledata)
				for i,st := range sheet.Rows[0].Cells{
					if st.Value == onlytitle { //首行为标题，单独拿出
						celid = i						//取出唯一索引对应的列数
					}else if st.Value == passtitle{
						pasid = i				////取出唯密码索引对应的列数
					}
				}
				userdata := sheet.Rows[userbind.Line].Cells
				//fmt.Println("已经绑定的行数：",userbind.Line)
				if userbind.Line > 0{
					linejc := sheet.Rows[userbind.Line].Cells[celid].Value
					fmt.Println("已存在的绑定：",linejc)
					if linejc == username{
						pasjc = sheet.Rows[userbind.Line].Cells[pasid].Value
						userdata = sheet.Rows[userbind.Line].Cells
						//fmt.Println(userdata)
					}else{
						fmt.Println("绑定关系出错开始查找并更新")
						for k, row := range sheet.Rows { //遍历工作表中的行
							if row.Cells[celid].Value == username{
								pasjc = row.Cells[pasid].Value
								userdata = row.Cells
								userbind.Binduser = username
								userbind.Line = k
								o.Update(&userbind)
								fmt.Println("找到数据第：",k,"行，数据内容：",userdata)
							}
						}
					}
				}else{
					for k, row := range sheet.Rows { //遍历工作表中的行
						if row.Cells[celid].Value == username{
							pasjc = row.Cells[pasid].Value
							userdata = row.Cells
							userbind.Binduser = username
							userbind.Line = k
							o.Insert(&userbind)
							fmt.Println("找到数据第：",k,"行，数据内容：",userdata)
						}
					}
				}
				//此处检查登陆成功逻辑
				if password == pasjc{
					max := len(titledata) - 1
					var looktitles string
					var lookdatas string
					datas := make(map[string]string)
					//fmt.Println("行数：",userbind.Line)
					if looktitle != ""{
						look := strings.Split(looktitle, ",")
						for _,v := range look{
							iv,_ := strconv.Atoi(v)
							if iv <= max{
								looktitles = looktitles + " " + titledata[iv].Value
								lookdatas = lookdatas + " " + userdata[iv].Value
								datas[titledata[iv].Value] = userdata[iv].Value
							}else{
								//fmt.Println("下标越界")
							}
						}
					}
					//fmt.Println(datas)
					mjson,_ :=json.Marshal(datas)
					mString :=string(mjson)
					fmt.Println("对应JSON：",mString)
					jsonre.Code = 200
					jsonre.Msg = "登陆成功，请稍等。。。"
					jsonre.Data = mString
				}else{
					jsonre.Msg = onlytitle + "或者登陆密码错误！"
				}
				//此处检查登陆成功逻辑
			}

		}
	}
	c.Data["json"] = &jsonre
	c.ServeJSON()
	//c.TplName = "Index/index.html"
}
/*
func (c *QQLogin) Get(){
	type Openidget struct {
		Openid string `json:"openid"`
	}
	type Userdata struct {
		Headimg string	`json:"figureurl_qq_2"`
		Nickname string	`json:"nickname"`
	}
	access_token := c.GetString("access_token")
	callback := "http://" + c.Ctx.Request.Host + "/qqlogin"
	api := "https://wauth.liuxiaoa.com/qqlogin.php?callback=" + callback
	if access_token == ""{
		c.Ctx.Redirect(302,api)
	}else{
		req := httplib.Get("https://graph.qq.com/oauth2.0/me?access_token=" + access_token)
		str,err := req.String()
		if err != nil{
			fmt.Println(err)
		}
		var myopenid Openidget
		openid := function.GetBetweenStr(str,"( "," )")
		json.Unmarshal([]byte(openid), &myopenid)
		req = httplib.Get("https://graph.qq.com/user/get_user_info?access_token=" + access_token + "&oauth_consumer_key=101576063&openid=" + myopenid.Openid)
		str,err = req.String()
		if err != nil{
			fmt.Println(err)
		}
		var userdata Userdata
		json.Unmarshal([]byte(str), &userdata)
		fmt.Println(userdata)
		c.Ctx.WriteString("获取成功")
	}

}*/