/**
@File:token_tools
@Author:LIUBIN
@Mail:liubin@wxjt.com.cn
@Time:19-8-18 下午2:03
@Software:GoLand
*/

//错误怠慢四、五位 02

/**
This file is a collection of tools for jwt functions
*/
package util

import (
	"auth/common"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"time"
)

//define secertkey of jwt
const (
	SECERTKEY = "sumgprinting"
)



//define TokenTools Class
type TokenTools struct {
	token       *jwt.Token
	tokenString string
	err         error
}

//function in TokenTools Class
//生成一个jwt标准的token
func (this *TokenTools) CreateToken(user string) (token string, err error) {
	//生成一个Token的指针
	this.token = jwt.New(jwt.SigningMethodHS256)
	//Token中claims信息
	claims := make(jwt.MapClaims)
	claims["username"] = user
	claims["exp"] = time.Now().Add(time.Hour*24 * time.Duration(1)).Unix()
	fmt.Println(time.Now().Add(time.Hour*24 * time.Duration(1)))
	claims["iat"] = time.Now().Unix()
	this.token.Claims = claims
	//加入安全码再次加密
	this.tokenString, this.err = this.token.SignedString([]byte(SECERTKEY))

	return this.tokenString, this.err
}

//function in TokenTools Class
//验证参数中传来的token是否正确
func (this *TokenTools) CheckToken(tokenString string) (*jwt.Token,bool, error) {
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjUwNzc2NzksImlhdCI6MTU2NTA3NDA3OSwidXNlcm5hbWUiOiIxMTEifQ.9wwXByWvuhZJ3vmMyJi6znbTjpkzIcxfUGSK_ltE__Q
	this.token, this.err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECERTKEY), nil
	})
	if this.err!=nil{
		return nil,false, fmt.Errorf("err of Parse jwt  |  %v",this.err)
	}

	return this.token, this.token.Valid, nil
}



func (this *TokenTools) CheckUserIdentity(tokenString string) (*common.UserIdentityResult,string,error) {

	UserIdentity := new(common.UserIdentityResult)

	jwtTokenRow, jwtTokenValid, checkTokenErr := this.CheckToken(tokenString)
	if checkTokenErr!=nil{
		CE(checkTokenErr)
		return nil,"",fmt.Errorf("check jwt-token err!  |  %v",checkTokenErr)
	}
	// 从原型中获取用户id
	username := jwtTokenRow.Claims.(jwt.MapClaims)["username"]

	// 连接redis
	redisAddr,err := GetSecertConf("reids::addr")
	if err!=nil{
		CE(err)
		return nil,"",fmt.Errorf("get redis addr err!  |  %v",err)
	}
	client := redis.NewClient(&redis.Options{Addr: redisAddr})
	// 查找redis中键为token的记录，返回记录的值
	redisRecode, readRedisErr := client.Get(tokenString).Result()
	if readRedisErr != nil {
		CE(readRedisErr)
		return nil,"",fmt.Errorf("read jwt-token from redis err!  |  %v",readRedisErr)
	}


	if redisRecode == username && jwtTokenValid {
		//todo:身份确认
		UserIdentity.ErrCode=1
		UserIdentity.ErrMsg = "Authentication passed !"
	} else {
		//todo:身份验证失败
		UserIdentity.ErrCode=-1
		UserIdentity.ErrMsg = "Authentication failed !"
	}

	{
		DP(redisRecode)
	}
	client.Close()

	return UserIdentity,redisRecode,nil


}
