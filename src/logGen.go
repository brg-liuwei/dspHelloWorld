package main

import (
	"encoding/json"
	"logger"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func currentTimeString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func uuid() string {
	return strconv.FormatInt(rand.Int63(), 10)
}

func price() string {
	return strconv.FormatInt(rand.Int63()%500, 10)
}

func genWinLog(l *logger.Log) {
	m := make(map[string]interface{})

	m["v"] = "1"
	m["log_type"] = "32"
	m["time_stamp"] = currentTimeString()
	m["ad_id"] = "cafefeed"
	m["order_id"] = "deadbeaf"
	m["exchange_user_id"] = "9"
	//m["dsp_user_id"] = ?
	m["media_type"] = "2"
	m["uuid"] = uuid()
	m["adexchange_id"] = "9" // yesky
	m["user_id"] = "abcd"
	m["user_agent"] = ""
	m["strike_price"] = strconv.Itoa(rand.Int() % 200)
	m["region_id"] = ""
	m["browser"] = "chrome"
	m["operation"] = "MacOS"
	m["language"] = "zh"
	m["agent_price"] = "110"
	m["advertiser_price"] = "120"
	m["reffer"] = ""
	m["adslot_id"] = strconv.Itoa(rand.Int() % 10)
	//m["adslot_position_relative"] = ?
	m["bid_id"] = uuid()
	m["price"] = m["strike_price"]
	m["key"] = "abcdef1234567"
	if b, e := json.Marshal(m); e != nil {
		return
	} else {
		l.Log(logger.INFO, string(b))
	}
}

func genBidLog(l *logger.Log) {
	m := make(map[string]interface{})

	m["v"] = "1"
	m["log_type"] = "30" // mango bid
	m["time_stamp"] = currentTimeString()
	m["uuid"] = uuid()
	m["bid_id"] = uuid()
	m["region_id"] = ""
	m["exchange_user_id"] = "8"
	m["user_agent"] = ""
	m["url"] = "http://www.dsp.com"
	m["language"] = "zh"
	m["media_type"] = "2"
	m["adslot_id"] = "100"
	m["adslot_size"] = "300X250"
	tc := make([]string, 0, 4)
	tc = append(tc, "1")
	tc = append(tc, "3")
	tc = append(tc, "5")
	m["target_channel"] = tc
	m["ad_id"] = "123"
	m["order_id"] = "456"
	m["creative_id"] = "789"
	m["creative_type_id"] = "111"
	m["vendor_type"] = "12"
	m["ad_class"] = make([]string, 0, 2)
	m["advertiser_name"] = "liuyu"
	m["creative_size"] = "300X250"
	m["creative_price"] = price()
	if b, e := json.Marshal(m); e != nil {
		return
	} else {
		l.Log(logger.INFO, string(b))
	}
}

func main() {
	if len(os.Args) != 3 {
		println("Usage: ", os.Args[0], " <bid log path> <win log path>")
		return
	}
	lw := logger.NewLog(os.Args[1])
	if lw == nil {
		panic("log path " + os.Args[1] + " error")
	}
	wlog := logger.NewLog(os.Args[2])
	if wlog == nil {
		panic("log path " + os.Args[2] + " error")
	}
	for {
		genBidLog(lw)
		genWinLog(wlog)
		time.Sleep(time.Second)
	}
}
