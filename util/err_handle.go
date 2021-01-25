package util

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"runtime"
	"strings"
)





func CE (err error) {

	_,file,line,_:=runtime.Caller(1)
	fileAbs:=strings.Split(file,"/")

	if (err!=nil) {
		fmt.Println(" ")
		fmt.Printf("\n%c[1;31m【ERR】 【%s】 %d  %c[0m",0x1B,fileAbs[len(fileAbs)-1],line, 0x1B)

		fmt.Printf("\n%c[1;31m", 0x1B)
		fmt.Println(err)
		fmt.Printf("%c[0m", 0x1B)
		logs.Error(err)
	}

}