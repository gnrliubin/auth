/*
File:		wx_tools

Author:		LIUBIN

Mail:		liubin@wxjt.com.cn

Time:		19-8-18 下午2:30

Software:	GoLand
*/

//错误怠慢四、五位 03

package util

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/config"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

//数据结构
//
// access_token api返回的json结构
//
//  test  a code block
//   fmt.print
//This is  a godoc test text
//
//This is a sigle-line model
//
//this is a description of struct AccessToken
//it has 4 data for the json struture what be reponsed from wx api server
//
// BUG(liubin):this is a bug test
type AccessToken struct {
	Errcode int64	`json:"errcode"`
	Token   string `json:"access_token"`
	Expires int64  `json:"expires_in"`
	Errmsg  string `json:"errmsg"`
}

//UserInfo 是微信返回用户信息的json结构
type UserInfo struct {
	ErrCode 	int8	`json:"errcode"`
	ErrMsg 		string 	`json:"errmsg"`
	UserId 		string 	`json:"UserId"`
	OpenId 		string 	`json:"OpenId"`
	DeviceId 	string 	`json:"DeviceId"`
}

//wx.config({
//	beta: true,// 必须这么写，否则wx.invoke调用形式的jsapi会有问题
//	debug: true, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
//	appId: '', // 必填，企业微信的corpID
//	timestamp: , // 必填，生成签名的时间戳
//	nonceStr: '', // 必填，生成签名的随机串
//	signature: '',// 必填，签名，见 附录-JS-SDK使用权限签名算法
//	jsApiList: [] // 必填，需要使用的JS接口列表，凡是要调用的接口都需要传进来
//});
type JsapiConfig struct {
	Beta 		bool	`json:"beta"`
	Debug 		bool	`json:"debug"`
	AppId 		string 	`json:"appId"`
	Timestamp 	int64	`json:"timestamp"`
	NonceStr 	uuid.UUID 	`json:"nonceStr"`
	Signature 	string	`json:"signature"`
}

// WxTools 是一个类
//
//wxapi工具类
//
//属性是向api发送请求时需要的一些参数
type WxTools struct {
	appid		string
	//redirectUri	string
	state 		string
	secret 		string
	token 		*AccessToken
}


//wxapi工具类初始化方法
//
//工厂模式创建类，从配置中读取企业微信信息
func NewWxTools() *WxTools{
	//读取配置
	configer, err := config.NewConfig("ini", "conf/secert.conf")
	if (err!=nil) {
		//fmt.Println("NewConfig of secert.....")
		//fmt.Println(err.Error())
		CE(err)
	}
	//初始化类及其属性
	// BUG(liubin): a bug in code
	tools := new(WxTools)
	tools.appid = configer.String("wx::appid")
	//tools.redirectUri = configer.String("wx::redirect_uri")
	tools.state  = configer.String("wx::state")
	tools.secret = configer.String("wx::secret")
	tools.token=new(AccessToken)
	//tools.getToken()

	// Output:
	//output test
	//返货类指针
	return tools
}

// 获取access_token
/**
 	token(and expire time) will saved in token.ini file ,after it be got
	when we want to use token and run this func ,it will check the token in token.ini file
	if expire time equal 0 or less than time now ,then func will require access_token from server
	or it will give you the data what you wanted from token.ini
 */
func (this *WxTools) getToken() (interface{},error)  {

	configer,err:= config.NewConfig("ini", "conf/token.ini")
	if (err!=nil) {
		CE(err)
		return nil,fmt.Errorf("fail to get access_token  ｜  %v",err)
	}
	//corp:="huirbv962J0rAfgLLFVDX3L_JpVOPDLfENbwZ8OZB70"
	corp:=this.secret
	//this.token=new(AccessToken)
	this.token.Token = configer.String("token::access_token")
	this.token.Expires, _ = configer.Int64("token::expires")

	//token过期时重新获取token
	if this.token.Expires == 0 || this.token.Expires <= time.Now().Unix() {
		uri := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", this.appid, corp)
		tokenByte := this.accessInterface("GET", uri)
		json.Unmarshal(tokenByte, &this.token)
		if this.token.Errcode > 0 {
			//请求错误处理
			return nil,fmt.Errorf("get access_token error! %q",this.token)
		}
		//token写入文件
		this.token.Expires = this.token.Expires + time.Now().Unix() - 1000
		configer.Set("token::access_token", this.token.Token)
		configer.Set("token::expires", fmt.Sprintf("%d", this.token.Expires))
		configer.SaveConfigFile("conf/token.ini")
	}
	{
		DP(this.token)
	}

	return this.token.Token,nil
}


//	Title			accessInterface
//	Description 	a http request
//	params 		method		string	method of request,a string with uppercase
//					uri			string	uri of api
//	return			respByte	[]byte	results of response from api serve
func (this *WxTools) accessInterface(method string, uri string) []byte{
	client := &http.Client{}
	request, _ := http.NewRequest(method,uri,nil)
	response, _ := client.Do(request)
	respByte ,e :=ioutil.ReadAll(response.Body)
	if(e!=nil){
		fmt.Println(e.Error())
	}
	return respByte
}

// CreateCodeUri 产生一个身份验证链接
func (this *WxTools) CreateCodeUri (redirectUri string ) string {

	uri:= fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect",this.appid,redirectUri,this.state)
	//uri:= fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect",this.appid,url.QueryEscape(this.redirectUri),this.state)

	return uri
}

//func (this *WxTools) GetUserInfo(code string) (*UserInfo,error) {
func (this *WxTools) GetUserInfo(code string) (map[string]interface{},error) {
	accessToken,err := this.getToken()
	if err!=nil{
		CE(err)
		return nil,fmt.Errorf("fail to get auth code  |  %v",err)
	}
	uri := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s",accessToken,code)
	userInfoByte := this.accessInterface("GET", uri)
	//userInfo := new(UserInfo)
	userInfo := make(map[string]interface{})

	json.Unmarshal(userInfoByte,&userInfo)

	{
		DP(userInfo)
	}
	if (userInfo["UserId"] == nil){
		return nil , fmt.Errorf("fail to get userid or openid  |  %v",userInfo)
	}

	return userInfo,nil
}

func (this *WxTools) GetUserDetail (code string) (map[string]interface{},error) {
	userDetail :=make(map[string]interface{})
	//userInfo := new(UserInfo)
	userInfo := make(map[string]interface{})
	userInfo,err := this.GetUserInfo(code)
	if err!=nil{
		CE(err)
		return nil ,fmt.Errorf("fail to get user detail  |  %v",err)
	}
	userId := userInfo["UserId"]

	accessToken,_ := this.getToken()
	uri := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s",accessToken,userId)

	userDetailByte := this.accessInterface("GET",uri)
	json.Unmarshal(userDetailByte,&userDetail)
	{
		DP(userDetail)
	}
	if (userDetail["userid"] == nil){
		return nil, fmt.Errorf("fail to get user detail  |  %v",userDetail)
	}

	return userDetail,nil

}


func (this *WxTools) GetJsapiTicket () (interface{},error) {

	var ticket = make(map[string]interface{})
	configer,err:= config.NewConfig("ini", "conf/token.ini")
	if (err!=nil) {
		CE(err)
		return nil,fmt.Errorf("fail to get jsapi_ticket  |  %v",err)
	}
	ticket["ticket"] = configer.String("ticket::ticket")
	ticket["expires"],_ = configer.Int64("ticket::expires")
	{
		DP(ticket)
	}

	// interface类型需要断言，否则与时间戳不能比较
	expires,_ := ticket["expires"].(int64)

	// 如果保存的ticket过期，则从api获取并存入文件
	if expires == 0 || expires <= time.Now().Unix() {
		// 获取token
		accessToken,err := this.getToken()
		if err != nil {
			return nil ,fmt.Errorf("get ticket error  |  %v",err)
		}
		// 请求API
		uri := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=%s",accessToken)
		ticketByte :=this.accessInterface("GET",uri)
		ticketFromApi := make(map[string]interface{})
		json.Unmarshal(ticketByte,&ticketFromApi)
		{
			DP(ticketFromApi)
		}
		// 返回的errcode存入ingerface类型，使用时需要断言
		// 如果错误代码不是0，做err处理
		errcode,_ := ticketFromApi["errcode"].(int)
		if errcode != 0{
			return nil ,fmt.Errorf("get ticket error  |  %v",ticketFromApi)
		}

		ticket["ticket"] = ticketFromApi["ticket"]
		// 传来的expire_in被存入interface后，使用时需要断言
		// 但是这个值断言成int64会报错，只能断言为float64
		expires_in,_ := ticketFromApi["expires_in"].(float64)
		// float64的值要与时间戳运算时，只能强行转换为int64
		ticket["expires"] = int64(expires_in) + time.Now().Unix() - 1000
		// 保存ticket
		configer.Set("ticket::ticket",ticket["ticket"].(string))
		configer.Set("ticket::expires",fmt.Sprintf("%d",ticket["expires"].(int64)))
		configer.SaveConfigFile("conf/token.ini")
	}

	return ticket["ticket"],nil
}

func (this *WxTools) GetJsapiSignature (noncestr uuid.UUID,uri string,timestamp int64) (string,error) {

	jsapiTicket,ticketErr := this.GetJsapiTicket()
	if ticketErr!=nil{
		CE(ticketErr)
		return "",fmt.Errorf("create jsapi signature error  |  %v",ticketErr)
	}
	signatureStr := fmt.Sprintf("jsapi_ticket=%s&noncestr=%v&timestamp=%d&url=%s",jsapiTicket,noncestr,timestamp,uri)
	fmt.Println(signatureStr)

	signatureSha := sha1.New()
	io.WriteString(signatureSha,signatureStr)
	// OUTPUT:&{[3186696664 1080495450 322735029 3720653562 75466895] [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 0 192}
	fmt.Println(signatureSha)
	signature := fmt.Sprintf("%x",signatureSha.Sum(nil))
	fmt.Println(signature)
	return signature,nil
}

func (this *WxTools) GetJsapiConfig (uri string) (*JsapiConfig,error) {
	jsapiConfig := new(JsapiConfig)
	noncestr := uuid.Must(uuid.NewV4())
	timestamp := time.Now().Unix()
	signature ,err := this.GetJsapiSignature(noncestr,uri,timestamp)
	if err!=nil{
		CE(err)
		return jsapiConfig,fmt.Errorf("create jsapi config err  |  %v",err)
	}

	jsapiConfig.Beta = true
	jsapiConfig.Debug = false
	jsapiConfig.AppId = this.appid
	jsapiConfig.Timestamp = timestamp
	jsapiConfig.NonceStr = noncestr
	jsapiConfig.Signature = signature
	{
		DP(jsapiConfig)
	}
	return jsapiConfig,nil

}