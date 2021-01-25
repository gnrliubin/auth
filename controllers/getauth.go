package controllers

import (
	"auth/common"
	"auth/models"
	"auth/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"

	//"strings"
)

type GetAuthController struct {
	beego.Controller
}

type RequestFrom struct {
	Paths []string		`json:"paths"`
 }

func (this *GetAuthController) Post() {

	results :=common.NewResults()
	//var rf RequestFrom
	//json.Unmarshal(this.Ctx.Input.RequestBody,&rf)
	//util.DP("get paths")
	//util.DP(rf)
	//验证用户省份
	jwtTools := new(util.TokenTools)
	// 获取请求头中的token
	token := this.Ctx.Input.Header("Authorization")
	jwtTokenString := strings.Split(token, " ")[1]
	// 得到token原型
	userIdentity,userid,err := jwtTools.CheckUserIdentity(jwtTokenString)
	if err!=nil {
		util.CE(err)
		results.ErrCode= 20301
		results.ErrMsg = fmt.Sprintf("check user identity err!  |  %v",err)
		this.Data["json"] = results
		this.ServeJSON()
	}
	results.Auth =userIdentity

	//读取权限列表
	o := orm.NewOrm()
	userAuth :=new([]*models.AuthGetUsersRules)

	_,oerr :=o.QueryTable("auth_get_users_rules").Filter("userid",userid).GroupBy("title").OrderBy("id").All(userAuth)
	if oerr!=nil{
		util.CE(oerr)
		results.ErrCode= 40302
		results.ErrMsg = fmt.Sprintf("check user identity err!  |  %v",oerr)
		this.Data["json"] = results
		this.ServeJSON()
	}


	results.Data = userAuth
	this.Data["json"] = results


	this.ServeJSON()
}
