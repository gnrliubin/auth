package controllers

import (

	"auth/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type TestController struct {
	beego.Controller
}


func (this *TestController) Post () {

	o := orm.NewOrm()
	userAuth :=new([]*models.AuthGetUsersRules)

	o.QueryTable("auth_get_users_rules").Filter("request_from","").All(userAuth)
	//o.Raw("select * from auth_get_users_rules").QueryRows(userAuth)
	fmt.Println(userAuth)
	this.Data["json"] = userAuth
	this.ServeJSON()
}