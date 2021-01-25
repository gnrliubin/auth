package util

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"runtime"
	"strings"
)

func DP(message interface{}) {

	_,file,line,_:=runtime.Caller(1)
	fileAbs:=strings.Split(file,"/")
	fmt.Println(" ")
	fmt.Printf("\n%c[1;32m【INF】 【%s】 %d  %c[0m",0x1B,fileAbs[len(fileAbs)-1],line, 0x1B)

	//fmt.Printf("\n%c[7;37m%s%c[0m\n\n", 0x1B, message, 0x1B)
	fmt.Printf("\n%c[1;37m", 0x1B)
	fmt.Println(message)
	fmt.Printf("%c[0m", 0x1B)
	logs.Info(message)
}
