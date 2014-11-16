package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	time.Sleep(time.Second)
	d := time.Since(t)
	fmt.Println("duration: ", int(d))
}
