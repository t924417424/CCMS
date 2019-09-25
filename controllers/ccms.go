package controllers

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"CCMS_O/function"
	"CCMS_O/models"
	"github.com/tealeg/xlsx"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/captcha"
)

var cpt *captcha.Captcha

type Clogin struct {
	Base
}
type Cindex struct {
	Base
}
type Left struct {
	Base
}
type Top struct {
	Base
}
type Main struct {
	Base
}
type Cupload struct {
	Base
}
type Cgetexcel struct {
	Base
}
type Sysset struct {
	Base
}
type Importdata struct {
	Base
}
type Datasys struct {
	Base
}
type Repass struct {
	Base
}
type Selects struct {
	Base
}
type Bind struct {
	Base
}
func (c *Clogin) Get(){
	c.TplName = "ccms/login.html"
}
func (c *Clogin) Post(){
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
		config := models.Config{}
		admin := models.Admin{}
		config.Id = 1
		admin.Id = 1
		admin.Username = username
		err := o.Read(&config)
		o.Read(&admin)
		if err == nil{
			password := function.Md5V(config.Salt + password)
			if password == admin.Password{
				jsonre.Code = 1
				jsonre.Msg = "登陆成功！"
				c.SetSession("suid",1)
			}else{
				jsonre.Msg = "账户或密码错误！"
			}
		}else{
			jsonre.Msg = "系统错误！请联系开发者！"
		}
	}
	c.Data["json"] = &jsonre
	c.ServeJSON()
}
func (c *Cindex) Get(){
	c.TplName = "ccms/index.html"
}
func (c *Left) Get(){
	c.TplName = "ccms/left.html"
}
func (c *Top) Get(){
	c.TplName = "ccms/top.html"
}
func (c *Main) Get(){
	c.Data["host"] = c.Ctx.Request.Host
	c.Data["var"] = "v1.0.1"
	c.Data["os"] = runtime.GOOS + " " + runtime.GOARCH
	c.Data["cpu"] = runtime.GOMAXPROCS(0)
	c.TplName = "ccms/main.html"
}
func (c *Cupload) Get(){
	c.TplName = "ccms/upload.html"
}
func (c *Cupload) Post(){
	type Jsons struct {
		Code int `json:"code"`
		Info string `json:"info"`
	}
	jsontest := &Jsons{200, "上传成功"}
	f, _, _ := c.GetFile("fileInfo") //获取上传的文件
	path := "./filetmp/" + "c48a93c545c647dc9c20a4420e6eee33"
	f.Close()                      //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	c.SaveToFile("fileInfo", path) //存文件
	text, _ := function.ReadAllIntoMemory(path)
	_, errs := xlsx.OpenBinary(text)
	key := []byte("scoresdcet111246141score")
	x1 := function.Encrypt3DES(text, key)
	//x2 := function.Decrypt3DES(x1, key)
	function.WriteWithIoutil(path, string(x1))
	_, err := os.Open(path)
	if err != nil {
		jsontest = &Jsons{100, "上传失败"}
	}else if errs != nil{
		jsontest = &Jsons{100, "文件格式错误"}
	}
	c.Data["json"] = jsontest
	c.ServeJSON()
}
func (c *Cgetexcel) Get(){
	var sheets []string
	path := "./filetmp/" + "c48a93c545c647dc9c20a4420e6eee33"
	text, _ := function.ReadAllIntoMemory(path)
	key := []byte("scoresdcet111246141score")
	x2 := function.Decrypt3DES(text, key)
	xlFile, err := xlsx.OpenBinary(x2)
	//fmt.Println(xlFile)
	if err != err{
		sheets = append(sheets, "文件读取失败，请重试！")
	}
	//var sheetname string
	i := 0
	for _, sheet := range xlFile.Sheets {
		sheets = append(sheets, sheet.Name)
		i++

	}
	c.Data["sheets"] = sheets
	c.TplName = "getsheet.html"
	c.TplName = "ccms/getexcel.html"
}
func (c *Sysset) Get(){
	c.TplName = "ccms/sys.html"
}
func (c *Sysset) Post(){
	jsonre := Newjson()
	sitename := c.GetString("sitename")
	if sitename != ""{
		o := orm.NewOrm()
		config := models.Config{}
		config.Id = 1
		config.Website = sitename
		_,err := o.Update(&config)
		if err == nil{
			jsonre.Code = 1
			jsonre.Msg = "修改成功！"
		}else{
			jsonre.Code = 0
			jsonre.Msg = "修改失败！"
		}
	}else{
		jsonre.Code = 0
		jsonre.Msg = "参数错误！"
	}
	c.Data["json"] = &jsonre
	c.ServeJSON()
}
func (c *Importdata) Post(){
	var title [][]*xlsx.Cell
	var titlesql string
	//var cells [][]*xlsx.Cell
	rowid,_ := c.GetInt("rowid")
	path := "./filetmp/" + "c48a93c545c647dc9c20a4420e6eee33"
	text, _ := function.ReadAllIntoMemory(path)
	key := []byte("scoresdcet111246141score")
	x2 := function.Decrypt3DES(text, key)
	xlFile, err := xlsx.OpenBinary(x2)
	sheet := xlFile.Sheets[rowid]	//选定Sheet工作表
	jsonre := Newjson()
	jsonre.Code = 0
	if err != nil{
		jsonre.Msg = "读取出错：" + err.Error()
	}else{
		//fmt.Println(sheet.Rows[0].Cells)
		//fmt.Println(typeof(sheet.Rows[0].Cells))
		for k, row := range sheet.Rows {	//遍历工作表中的行
			if k == 0{						//首行为标题，单独拿出
				title = append(title,row.Cells)
				for k,v := range row.Cells{
					if k < len(row.Cells) && k > 0{
						titlesql = titlesql + ","
					}
					titlesql = titlesql + v.String()
				}
			}
			continue
		}
		o := orm.NewOrm()
		title := models.Dataset{}
		title.Id = 1
		o.Read(&title)
		title.Exceltitle = titlesql
		_,err := o.Update(&title)
		if err == nil{
			jsonre.Code = 1
			jsonre.Msg = "导入成功"
		}else{
			jsonre.Code = 0
			jsonre.Msg = "导入失败：" + err.Error()
		}
	}
	/*c.Data["cells"] = cells
	c.TplName = "getcell.html"*/
	c.Data["json"] = &jsonre
	c.ServeJSON()
}
func (c *Datasys) Get(){
	o := orm.NewOrm()
	excel := models.Dataset{}
	excel.Id = 1
	o.Read(&excel)
	looktitle := ""
	only := 0
	pass := 0
	title := strings.Split(excel.Exceltitle, ",")
	max := len(title) - 1
	if excel.Looktitle != ""{
		look := strings.Split(excel.Looktitle, ",")
		for _,v := range look{
			iv,_ := strconv.Atoi(v)
			if iv <= max{
				looktitle = looktitle + " " + title[iv]
			}else{
				looktitle = looktitle + " " + "越界的下标" + strconv.Itoa(iv)
			}
		}
	}
	if excel.Onlytitle != ""{
		for k,v := range title{
			if v == excel.Onlytitle{
				only = k
			}else if v == excel.Passtitle{
				pass = k
			}
		}
	}
	c.Data["look"] = looktitle
	c.Data["only"] = only
	c.Data["pass"] = pass
	c.Data["excel"] = title
	c.Data["excels"] = &excel
	c.TplName = "ccms/datasys.html"
}
func (c *Datasys) Post(){
	look := c.GetString("etitle")
	onlytitle := c.GetString("onlydata")
	onlypass := c.GetString("onlypass")
	jsonre := Newjson()
	if look != "" && onlytitle != "" && onlypass != ""{
		o := orm.NewOrm()
		excel := models.Dataset{}
		excel.Id = 1
		o.Read(&excel)
		title := strings.Split(look, ",")
		titles := strings.Split(excel.Exceltitle, ",")
		maxlen := len(titles) - 1
		status := 1
		for _,v := range title{
			iv,err := strconv.Atoi(v)
			if err == nil{
				if iv > maxlen{
					status = 0
					jsonre.Code = 0
					jsonre.Msg = "展示标题行下标越界！"
					break
				}
			}else{
				status = 0
				jsonre.Code = 0
				jsonre.Msg = "展示标题行下标必须为Int类型！" + err.Error()
				break
			}
		}
		for _,v := range titles{
			if v != onlytitle{
				status = 0
			}else{
				status = 1
				break
			}
		}
		for _,v := range titles{
			if v != onlypass{
				status = 0
			}else{
				status = 1
				break
			}
		}
		/*if err == nil{
			if ov > maxlen{
				status = 0
				jsonre.Code = 0
				jsonre.Msg = "唯一不重复标题行下标设置越界！"
			}
		}else{
			status = 0
			jsonre.Code = 0
			jsonre.Msg = "展示标题行下标必须为Int类型！"
		}*/
		if status == 1{
			excel.Looktitle = look
			excel.Onlytitle = onlytitle
			excel.Passtitle = onlypass
			if _,err := o.Update(&excel); err == nil{
				jsonre.Code = 1
				jsonre.Msg = "更新成功！"
			}else{
				jsonre.Code = 0
				jsonre.Msg = "更新失败：" + err.Error()
			}
		}else{
			jsonre.Code = 0
			jsonre.Msg = "唯一列" + onlytitle + "或者密码列" + onlypass + "不存在当前标题行中！"
		}
	}else{
		jsonre.Code = 0
		jsonre.Msg = "设置内容不能为空！"
	}
	c.Data["json"] = &jsonre
	c.ServeJSON()
}
func (c *Repass) Get(){
	c.TplName = "ccms/repass.html"
}
func (c *Repass) Post(){
	oldpass := c.GetString("oldpass")
	newpass := c.GetString("newpass")
	repass := c.GetString("repass")
	jsonre := Newjson()
	jsonre.Code = 0
	if oldpass == ""||newpass == ""||repass == ""{
		jsonre.Msg = "内容不能为空！"
	}else if newpass != repass{
		jsonre.Msg = "两次密码输入不一致！"
	}else{
		o := orm.NewOrm()
		config := models.Config{}
		admin := models.Admin{}
		config.Id = 1
		admin.Id = 1
		o.Read(&config)
		o.Read(&admin)
		password := function.Md5V(config.Salt + oldpass)
		if password == admin.Password{
			admin.Password = function.Md5V(config.Salt + newpass)
			o.Update(&admin)
			jsonre.Code = 1
			jsonre.Msg = "密码修改成功！"
		}else{
			jsonre.Msg = "当前密码验证不通过！"
		}
	}
	c.Data["json"] = &jsonre
	c.ServeJSON()
}
func (c *Selects) Get(){
	o := orm.NewOrm()
	set := models.Dataset{}
	set.Id = 1
	o.Read(&set)
	c.Data["only"] = set.Onlytitle
	c.TplName = "ccms/select.html"
}
func (c *Selects) Select() {
	sel := c.GetString(":select")
	if sel == ""{
		c.Ctx.WriteString("检索内容不能为空！")
		return
	}
	o := orm.NewOrm()
	dataset := models.Dataset{}
	dataset.Id = 1
	o.Read(&dataset)
	if dataset.Onlytitle == "" {
		c.Ctx.WriteString("请先设置唯一不重复字段！")
		return
	}
	var title [][]*xlsx.Cell
	var cells [][]*xlsx.Cell
	path := "./filetmp/" + "c48a93c545c647dc9c20a4420e6eee33"
	text, _ := function.ReadAllIntoMemory(path)
	key := []byte("scoresdcet111246141score")
	x2 := function.Decrypt3DES(text, key)
	xlFile, err := xlsx.OpenBinary(x2)
	stid := 0
	celid := 0
	status := 0
	for s, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows { //遍历工作表中的行
			if k == 0 {
				for i,st := range row.Cells{
					if status == 1{
						continue
					}
					if st.Value == dataset.Onlytitle { //首行为标题，单独拿出
						stid = s
						celid = i
						status = 1
						continue
					}
				}
			}
		}
	}
	//以上步骤为遍历唯一索引存在的工作表的ID和表中的指定列
	//i := 0
	//rowid := 0
	sheet := xlFile.Sheets[stid]	//选定Sheet工作表

	if err != nil {
		c.Ctx.WriteString("读取出错")
	} else {
		//fmt.Println(sheet.Rows[0].Cells)
		//fmt.Println(typeof(sheet.Rows[0].Cells))
		for k, row := range sheet.Rows { //遍历工作表中的行
			if k == 0 { //首行为标题，单独拿出
				title = append(title, row.Cells)
				//fmt.Println(len(row.Cells))
			} else { //遍历工作表的每一列
				if row.Cells[celid].Value == sel{
					cells = append(cells, row.Cells)
				}
			}
		}
	}
	c.Data["cells"] = title
	c.Data["cellss"] = cells
	c.TplName = "ccms/getcell.html"
	//fmt.Println(cells)
}
func (c *Bind) Get(){
	o := orm.NewOrm()
	set := models.Dataset{}
	set.Id = 1
	o.Read(&set)
	c.Data["only"] = set.Onlytitle
	c.TplName = "ccms/bind.html"
}
func (c *Bind) Post(){
	o := orm.NewOrm()
	_,err := o.Raw("truncate bind_uid").Exec()
	jsonre := Newjson()
	jsonre.Code = 0
	if err == nil{
		jsonre.Code = 1
		jsonre.Msg = "清理完毕"
	}else{
		jsonre.Msg = "出现错误： " + err.Error()
	}
	c.Data["json"] = &jsonre
	c.ServeJSON()
}
func (c *Bind) Bind(){
	sel := c.GetString(":select")
	if sel == ""{
		c.Ctx.WriteString("检索内容不能为空！")
		return
	}
	o := orm.NewOrm()
	bind := models.BindUid{}
	sys := models.Dataset{}
	sys.Id = 1
	o.Read(&sys)
	bind.Binduser = sel
	o.Read(&bind,"Binduser")
	c.Data["binddata"] = sys.Onlytitle
	c.Data["bind"] = bind
	c.TplName = "ccms/getbind.html"
}