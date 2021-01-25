package controllers

import (
	"auth/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"time"
)

type WxLoginController struct {
	beego.Controller
}



func (this *WxLoginController) Post() {
	results := map[string]interface{}{"errCode":0,"errMsg":"ok","data":nil}
	//results := new(Results)
	//results.ErrCode = 0
	//results.ErrMsg = "ok"
	//results.Data = make(map[string]interface{})

	code := this.GetString("code")

	//userInfo := new(util.UserInfo)
	userDetail := make(map[string]interface{})
	var wxTools = util.NewWxTools()
	//userInfo = wxTools.GetUserDetail(code)
	userDetail, err := wxTools.GetUserDetail(code)
	if err != nil {
		util.CE(err)
		results["errCode"] = 20101
		results["errMsg"] = fmt.Errorf("get user detail err!  |  %v",err).Error()
		//results.ErrCode = 20103
		//results.ErrMsg = err.Error()
		this.Data["json"] = results
		this.ServeJSON()
	}

	// 生成jwt-token
	jwtTools := &util.TokenTools{}
	jwtToken, jwtErr := jwtTools.CreateToken(userDetail["userid"].(string))

	if jwtErr!= nil{
		results["errCode"] = 20102
		results["errMsg"] = fmt.Errorf("create jwt token err!  |  %v",jwtErr).Error()
		this.Data["json"] = results
		this.ServeJSON()
	}

	{
		util.DP(userDetail["userid"].(string))
	}

	// 设置过期时间
	//now := time.Now().Unix()
	var expire time.Duration = 7200
	userDetail["jwt"] = jwtToken
	userDetail["expire"] = expire
	fmt.Println(code)
	results["data"] = userDetail
	//results.Data = userDetail
	this.Data["json"] = results

	redisAddr, err := util.GetSecertConf("redis::addr")
	if err != nil {
		results["errCode"] = 20103
		results["errMsg"] = err.Error()
		//results.ErrCode = 20101
		//results.ErrMsg = err.Error()
		this.Data["json"] = results
		this.ServeJSON()
	}
	client := redis.NewClient(&redis.Options{Addr: redisAddr})

	nx := client.Set(jwtToken, userDetail["userid"], expire*time.Second).Err()
	if nx != nil {
		util.CE(fmt.Errorf("save jwt to redis error  |  %v", nx))
		results["errCode"] = 30101
		results["errMsg"] = fmt.Errorf("set redis err!  |  %v",nx).Error()
		//results.ErrCode = 10101
		//results.ErrMsg = nx.Error()
		this.Data["json"] = results
		this.ServeJSON()
	}
	defer client.Close()

	fmt.Println(this.Data)
	//this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin","*")
	this.ServeJSON()
}
