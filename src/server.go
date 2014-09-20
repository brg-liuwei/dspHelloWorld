package main

import (
	"logger"
	"manager"

	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if ok := logger.Init("dspLog.log"); !ok {
		panic("logger init error.")
	}

	// test manager
	c := manager.NewCommand()
	jstr := `{
        "oper_type":"4",
        "fmt_ver":"1",
        "data":
        [   
        "adfaf2131dfafb",
        "adfaf2131dfafb",
        "adfaf2131dfafb"
        ]   
    }`
	if c.Parse(jstr) {
		fmt.Println("order_type: ", c.Ctype)
		fmt.Println("version: ", c.Cversion)
		fmt.Println("data: ", c.Data)
	} else {
		fmt.Println("Parse error")
	}
	// end of test manager
}
