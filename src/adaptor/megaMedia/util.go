package megaMedia

import (
	"aes"
	"bid"
	"logger"

	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ErrorRep(w http.ResponseWriter) {
	errMsg := "Illegal Request"
	w.Header().Set("Content-Length", fmt.Sprintf("%d", errMsg))
	w.Write([]byte(errMsg))
}

func OkRep(w http.ResponseWriter) {
	errMsg := "ok"
	w.Header().Set("Content-Length", fmt.Sprintf("%d", errMsg))
	w.Write([]byte(errMsg))
}

func YeskyClickHandler(w http.ResponseWriter, r *http.Request) {
	megaMediaLogger.Log(logger.INFO, "[click] mega: ", *r)
	OkRep(w)
}

func YeskyDisplayHandler(w http.ResponseWriter, r *http.Request) {
	megaMediaLogger.Log(logger.INFO, "[display] mega", *r)
	OkRep(w)
}

func YeskyBidHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	//fmt.Println(t, "yesky request: ", *r)
	bidRequest := NewBidRequest(r)
	if bidRequest == nil {
		goto giveup
	}
	//fmt.Printf("\nyesky bidrequest: %#v\n", *bidRequest)

	commonRequest := bidRequest.ParseToCommon()
	//fmt.Printf("\ncommon bidrequest: %#v\n", *commonRequest)

	commonResponse, isBid := bid.Bid(commonRequest)
	//fmt.Printf("\ncommon response: %#v\n", *commonResponse)

	if !isBid {
		goto giveup
	}

	bidResponse := new(MgxBidResponse)
	bidResponse.ParseFromCommon(commonResponse)

	bidResponse.Response(w)
	thinkTime := time.Since(t)
	//fmt.Println("thinking time: ", thinkTime)
	m := make(map[string]interface{})
	m["v"] = 1
	m["time_stamp"] = logger.CurrentTimeString()
	m["exchange_user_id"] = strconv.Itoa(int(common.MEGAMEDIA))
	m["think_time"] = strconv.Itoa(int(thinkTime) / 1000000) /* ms */

bidlog:
	m["log_type"] = "30"
	m["uuid"] = logger.UUID()
	m["bid_id"] = commonResponse.BidId
	m["ad_id"] = commonResponse.AdId
	m["order_id"] = commonResponse.OrderId
	m["creative_price"] = commonResponse.Ads[0].Price
	logger.BidLog.JsonLog(logger.INFO, m)
	return

giveuplog:
	m["log_type"] = "31"
	m["uuid"] = ""
	logger.GiveupLog.JsonLog(logger.INFO, m)
}

func YeskyWinHandler(w http.ResponseWriter, r *http.Request) {
	var price int
	var eprice string
	var ext, adId, orderId string

	if r.Method != "GET" {
		megaMediaLogger.Log(logger.ERROR, "mega win notice method err: ", r.Method, "; r: ", *r)
		goto end_error
	}
	if err := r.ParseForm(); err != nil {
		megaMediaLogger.Log(logger.ERROR, "mega win notice parse form err: ", err)
		goto end_error
	}
	if ext = r.FormValue("extend_data"); ext != "yesky" {
		megaMediaLogger.Log(logger.ERROR, "mega win notice ext illegal:", ext)
		goto end_error
	}
	adId = r.FormValue("ad_id")
	orderId = r.FormValue("order_id")
	eprice = r.FormValue("win_price")
	if len(strings.TrimSpace(eprice)) == 0 {
		megaMediaLogger.Log(logger.ERROR, "mega win notice parse form no win price")
		goto end_error
	}
	price = aes.GetDecryptedPrice(eprice)

	megaMediaLogger.Log(logger.INFO, "mega win price: ", price, "adId:", adId, "orderId: ", orderId)

	m := make(map[string]interface{})
	m["v"] = "1"
	m["log_type"] = "32"
	m["time_stamp"] = logger.CurrentTimeString()
	m["ad_id"] = adId
	m["order_id"] = orderId
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
