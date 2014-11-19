package megaMedia

import (
	//"aes"
	"bid"
	"common"
	"logger"

	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ErrorRep(w http.ResponseWriter) {
	errMsg := "Illegal Request"
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(errMsg)))
	w.Write([]byte(errMsg))
}

func OkRep(w http.ResponseWriter) {
	errMsg := "ok"
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(errMsg)))
	w.Write([]byte(errMsg))
}

func YeskyClickHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	if r.Method != "GET" {
		megaMediaLogger.Log(logger.ERROR, "mega click notice method err: ", r.Method)
		ErrorRep(w)
		return
	}
	if err := r.ParseForm(); err != nil {
		megaMediaLogger.Log(logger.ERROR, "mega click notice parse form err: ", err)
		ErrorRep(w)
		return
	}

	adId := r.FormValue("ad_id")
	orderId := r.FormValue("order_id")
	bidId := r.FormValue("bid_id")

	m["v"] = "1"
	m["log_type"] = "9"
	m["time_stamp"] = logger.CurrentTimeString()
	m["ad_id"] = adId
	m["order_id"] = orderId
	m["exchange_user_id"] = strconv.Itoa(int(common.MEGAMEDIA))
	m["uuid"] = logger.UUID()
	m["bid_id"] = bidId

	ip := r.Header.Get("X-FORWARDED-FOR") /* this is real ip */
	if len(ip) == 0 {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	m["user_ip"] = ip

	logger.DisplayLog.JsonLog(logger.INFO, m)
	OkRep(w)
}

func YeskyDisplayHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	if r.Method != "GET" {
		megaMediaLogger.Log(logger.ERROR, "mega display notice method err: ", r.Method)
		ErrorRep(w)
		return
	}
	if err := r.ParseForm(); err != nil {
		megaMediaLogger.Log(logger.ERROR, "mega display notice parse form err: ", err)
		ErrorRep(w)
		return
	}
	adId := r.FormValue("ad_id")
	orderId := r.FormValue("order_id")
	bidId := r.FormValue("bid_id")

	m["v"] = "1"
	m["log_type"] = "8"
	m["time_stamp"] = logger.CurrentTimeString()
	m["ad_id"] = adId
	m["order_id"] = orderId
	m["exchange_user_id"] = strconv.Itoa(int(common.MEGAMEDIA))
	m["uuid"] = logger.UUID()
	m["bid_id"] = bidId

	ip := r.Header.Get("X-FORWARDED-FOR") /* this is real ip */
	if len(ip) == 0 {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	m["user_ip"] = ip

	logger.DisplayLog.JsonLog(logger.INFO, m)
	OkRep(w)
}

func YeskyBidHandler(w http.ResponseWriter, r *http.Request) {
	var commonRequest *common.BidRequest
	var commonResponse *common.BidResponse
	var isBid bool
	var thinkTime time.Duration

	//fmt.Println(t, "yesky request: ", *r)

	t := time.Now()
	bidResponse := new(MgxBidResponse)
	m := make(map[string]interface{})

	bidRequest := NewBidRequest(r)
	if bidRequest == nil {
		goto giveup
	}
	//fmt.Printf("\nyesky bidrequest: %#v\n", *bidRequest)

	commonRequest = bidRequest.ParseToCommon()
	//fmt.Printf("\ncommon bidrequest: %#v\n", *commonRequest)

	commonResponse, isBid = bid.Bid(commonRequest)
	//fmt.Printf("\ncommon response: %#v\n", *commonResponse)

	if !isBid || len(commonResponse.Ads) == 0 {
		goto giveup
	}

	bidResponse.ParseFromCommon(commonResponse)

	bidResponse.Response(w)
	thinkTime = time.Since(t)
	//fmt.Println("thinking time: ", thinkTime)
	m["v"] = 1
	m["time_stamp"] = logger.CurrentTimeString()
	m["exchange_user_id"] = strconv.Itoa(int(common.MEGAMEDIA))
	m["think_time"] = strconv.Itoa(int(thinkTime) / 1000000) /* ms */

	//bidlog:
	m["log_type"] = "30"
	m["uuid"] = logger.UUID()
	m["bid_id"] = commonResponse.BidId
	m["ad_id"] = commonResponse.Ads[0].AdId
	m["order_id"] = commonResponse.Ads[0].OrderId
	m["creative_price"] = commonResponse.Ads[0].Price
	logger.BidLog.JsonLog(logger.INFO, m)
	return

giveup:
	m["log_type"] = "31"
	m["uuid"] = ""
	logger.GiveupLog.JsonLog(logger.INFO, m)
}

func YeskyWinHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var price int
	var eprice string
	var ext, adId, orderId, bidId string

	m := make(map[string]interface{})

	fmt.Println("Win handler: r: ", *r)

	if r.Method != "GET" {
		megaMediaLogger.Log(logger.ERROR, "mega win notice method err: ", r.Method, "; r: ", *r)
		goto end_error
	}
	if err = r.ParseForm(); err != nil {
		megaMediaLogger.Log(logger.ERROR, "mega win notice parse form err: ", err)
		goto end_error
	}
	adId = r.FormValue("ad_id")
	orderId = r.FormValue("order_id")
	bidId = r.FormValue("bid_id")
	eprice = r.FormValue("win_price")
	if len(eprice) == 0 || len(strings.TrimSpace(eprice)) == 0 {
		megaMediaLogger.Log(logger.ERROR, "mega win notice parse form no win price")
		goto end_error
	}
	if ext = r.FormValue("extend_data"); ext != "yesky" {
		megaMediaLogger.Log(logger.ERROR, "mega win notice ext illegal:", ext)
		goto end_error
	}
	fmt.Println("before decry: ", eprice, " eprice len = ", len(eprice))
	//price = aes.GetDecryptedPrice(eprice)
	price = getDecryptedPriceFromServer(eprice)
	fmt.Println("--------> after decry price: ", price)

	megaMediaLogger.Log(logger.INFO, "mega win price: ", price, "adId:", adId, "orderId: ", orderId)

	m["v"] = "1"
	m["log_type"] = "32"
	m["time_stamp"] = logger.CurrentTimeString()
	m["ad_id"] = adId
	m["order_id"] = orderId
	m["bid_id"] = bidId
	m["exchange_user_id"] = strconv.Itoa(int(common.MEGAMEDIA))
	m["dsp_user_id"] = ""
	m["media_type"] = ""
	m["uuid"] = logger.UUID()
	m["price"] = strconv.Itoa(price)
	logger.WinLog.JsonLog(logger.INFO, m)

	OkRep(w)
	return

end_error:
	ErrorRep(w)
	return
}

func Reader(r io.Reader, ch chan int) {
	var buf [1024]byte
	n, err := r.Read(buf[:])
	if err != nil {
		fmt.Println("in Reader err: ", err)
		ch <- 0
		return
	}
	if n >= 1024 {
		n = 1024
	}
	fmt.Println("Read data: ", string(buf[:n]))
	rc, _ := strconv.Atoi(string(buf[:n]))
	ch <- rc
	return
}

func getDecryptedPriceFromServer(eprice string) int {
	c, err := net.Dial("unix", "/tmp/unix_aes_listen")
	if err != nil {
		fmt.Println("decrypt conn error: ", err)
		return 0
	}
	defer c.Close()

	ch := make(chan int)
	go Reader(c, ch)

	if _, err = c.Write([]byte(eprice)); err != nil {
		fmt.Println("decrypt write error: ", err)
		return 0
	}
	return <-ch
}
