package common

//用户身份验证结果结构体
type UserIdentityResult struct {
	ErrCode int64  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}



//展示结果的结构体
//使用map[string]interface{}{"errCode":0,"errMsg":"ok","data":nil}也可以
//结果json结构如下
//{
//	errCode:0
//	errMsg:"ok"
//	auth:{
//		errCode:0			//0：不需要身份验证；1：身份验证通过；-1：身份验证失败
//		errMsg:"no need Authentication"
//	}
//	data:{
//		各种请求返回的结果
//	}
//}

type Results struct {
	ErrCode int64               `json:"errCode"`
	ErrMsg  string              `json:"errMsg"`
	Auth    *UserIdentityResult `json:"auth"`
	Data    interface{}         `json:"data"`
}

func NewResults() *Results {
	results := new(Results)
	auth := new(UserIdentityResult)
	auth.ErrCode = 0
	auth.ErrMsg = "no need Authentication"
	results.ErrCode = 0
	results.ErrMsg = "ok"
	results.Auth = auth
	results.Data = nil
	return results
}
