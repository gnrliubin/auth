package models

import "github.com/astaxie/beego/orm"

type AuthGetUsersRules struct {
	Id    uint   `orm:"column(id)"`		//一定要有个id字段
	Name         string `orm:"column(name);size(50)" description:"姓名"`
	Userid       string `orm:"column(userid);size(50)" description:"企业微信id"`
	RequestFrom  string `orm:"column(request_from);size(255)" description:"路由，当前页面"`
	Title        string `orm:"column(title);size(50)"`
	RouteName    string `orm:"column(route_name);size(100);null" description:"menu路由名称"`
	Component    string `orm:"column(component);size(100);null" description:"menu路由需要的组件"`
	RuleTitle    string `orm:"column(rule_title);size(50)"`
	Path         string `orm:"column(path);size(100)" description:"menu连接的路由；button后端地址"`
	Type         uint8  `orm:"column(type)" description:"0 menu；1 button"`
	Icon         string `orm:"column(icon);size(100)" description:"图标"`
	Platform     uint   `orm:"column(platform)" description:"客户端平台"`
	PlatformName string `orm:"column(platform_name);size(50)" description:"名称"`
}

func (this *AuthGetUsersRules) TableName() string{
	return "auth_get_users_rules"
}

func init() {
	orm.RegisterModel(new(AuthGetUsersRules))
}