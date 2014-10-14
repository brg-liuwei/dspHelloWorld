package main

import (
	"fmt"
	"net/http"
	"strings"
)

// func main() {
// 	reader := strings.NewReader("key=encoding&val=x-www-form-urlencoded")
// 	//response, err := http.Post("http://kv-store:8888/levelput", "application/octet-stream", reader)
// 	//response, err := http.Post("http://127.0.0.1:8888/levelput", "application/postscript", reader)
// 	response, err := http.Post("http://127.0.0.1:8888/levelput", "application/x-www-form-urlencoded", reader)
// 	if err != nil {
// 		fmt.Println("submit error:", err)
// 	} else {
// 		fmt.Println(response)
// 	}
// }

var bidRequest string = `{
    "id": "d6e608bcbe5344d4b5135f710ed65fb9", "imp": [
    {
        "impid": "483ceb489810452187f72667509b7b0d", "bidfloor": 6000,
        "bidfloorcur": "CNY",
        "w": 640,
        "h": 100,
        "instl": 0
    } ],
    "app": {
        "aid": "9eded84882b9441583f8f623aff2fe8e",
        "name": "App Name",
        "cat": [ "2" ],
        "ver": "1.0",
        "bundle": "com.company.appname",
        "paid": 0,
        "storeurl": "https://itunes.apple.com/cn/app/id12345678?mt=8"
    }, "device": {
        "did": "75184389584b32a37f4f4570ca85112da1463707",
        "dpid": "e1ef08c816e8e1604c7b4c5ddad8cdaf2edfc843",
        "ua": "Mozilla%2F5.0+%28iPhone%3B+U%3B+CPU+like+Mac+OS+X%3B +en%29+AppleWebKit%2F420% 2B+%28KHTML",
        "ip": "123.123.123.123", "country": "CN", "carrier": "46000", "language": "zh", "make": "Apple", "model": "iPhone5,1", "os": "iOS",
        "osv": "7.0", "connectiontype": 2,
        "devicetype": 1,
        "loc": "38.04165,114.50884" }
    }`

//var postAddr string = "http://localhost:12345/posttest"
//var postAddr string = "http://localhost:18124/mango"
var postAddr string = "http://124.232.133.211:18124/mango/bid"

func main() {
	reader := strings.NewReader(bidRequest)
	fmt.Println("bidRequest len: ", len(bidRequest))
	r, e := http.Post(postAddr, "application/json", reader)
	if e != nil {
		panic(e)
	}
	defer r.Body.Close()
}
