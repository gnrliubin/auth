
//错误怠慢四、五位 01


package util

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

func GetSecertConf(key string) (string, error) {

	configer, err := config.NewConfig("ini", "conf/secert.conf")
	if err != nil {
		CE(err)
		return "", fmt.Errorf("get secert conf err  |  %v", err)
	}
	results := configer.String("redis::addr")
	return results, nil
}
