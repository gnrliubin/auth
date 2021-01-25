// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"auth/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/weixin",
			beego.NSRouter("/login",&controllers.WxLoginController{}),  // 错误代码二、三位 01
			beego.NSRouter("/jsapi",&controllers.JsapiController{}),	// 错误代码二、三位 02
			beego.NSRouter("/getauth",&controllers.GetAuthController{}),// 错误代码二、三位 03
		),
		//beego.NSNamespace("/account",
		//	//beego.NSInclude(
		//	//	&controllers.UserController{},
		//	//),
		//),
		beego.NSRouter("/test",&controllers.TestController{}),
	)
	beego.AddNamespace(ns)
}
