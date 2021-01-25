package controllers

import (
	"auth/common"
	"auth/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)



type JsapiController struct {
	beego.Controller
}


func (this *JsapiController) Post() {
	//定义返回结构
	//var results common.Results
	//results.ErrCode= 0
	//results.ErrMsg = "ok"
	//results.Data = nil

	results :=common.NewResults()

	//获取请求中的参数
	var requestData map[string]string
	ReBodyErr := json.Unmarshal(this.Ctx.Input.RequestBody, &requestData)
	if ReBodyErr != nil {
		util.CE(ReBodyErr)
		results.ErrCode = 10201
		results.ErrMsg = fmt.Errorf("get request params err!  |  %v",ReBodyErr).Error()
		this.Data["json"] = results
		this.ServeJSON()
	}
	uri := requestData["uri"]

	{
		util.DP(requestData)
		util.DP(uri)
	}

	wxTools := util.NewWxTools()

	configStruct, err := wxTools.GetJsapiConfig(uri)
	if err != nil {
		util.CE(err)
		results.ErrCode = 20202
		results.ErrMsg = fmt.Errorf("get jsapi config err!  |  %v",err).Error()
		this.Data["json"] = results
		this.ServeJSON()
	}

	// 转json字符串的话，前端需要再转成对象，用map或者struct可以直接转为json对象给前端
	//configBytes, structToBytesErr := json.Marshal(configStruct)
	//if structToBytesErr != nil {
	//	util.CE(structToBytesErr)
	//}
	//{
	//	util.DP(string(configBytes))
	//}
	//this.Data["json"] = string(configBytes)

	//直接传结构体，前端直接接收json对象
	results.Data = configStruct
	this.Data["json"] = results
	{
		util.DP(this.Data["json"])
	}

	this.ServeJSON()

}
