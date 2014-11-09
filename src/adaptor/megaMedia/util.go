package megaMedia

import (
	"aes"
	"bid"
	"logger"

	"fmt"
	"net/http"
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
	megaMediaLogger.Log(logger.INFO, "[click] mega")
	OkRep(w)
}

func YeskyDisplayHandler(w http.ResponseWriter, r *http.Request) {
	megaMediaLogger.Log(logger.INFO, "[display] mega")
	OkRep(w)
}

func YeskyBidHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t, "yesky request: ", *r)
	bidRequest := NewBidRequest(r)
	if bidRequest == nil {
		return
	}
	fmt.Printf("\nyesky bidrequest: %#v\n", *bidRequest)

	commonRequest := bidRequest.ParseToCommon()
	fmt.Printf("\ncommon bidrequest: %#v\n", *commonRequest)

	commonResponse := bid.Bid(commonRequest)
	fmt.Printf("\ncommon response: %#v\n", *commonResponse)

	bidResponse := new(MgxBidResponse)
	bidResponse.ParseFromCommon(commonResponse)

	bidResponse.Response(w)
	fmt.Println("======> delta time: ", time.Since(t))
}

func YeskyWinHandler(w http.ResponseWriter, r *http.Request) {
	var price int
	var eprice string

	if r.Method != "GET" {
		megaMediaLogger.Log(logger.ERROR, "mega win notice method err: ", r.Method)
		goto end_error
	}
	if err := r.ParseForm(); err != nil {
		megaMediaLogger.Log(logger.ERROR, "mega win notice parse form err: ", err)
		goto end_error
	}
	eprice = r.FormValue("win_price")
	if len(strings.TrimSpace(eprice)) == 0 {
		megaMediaLogger.Log(logger.ERROR, "mega win notice parse form no win price")
		goto end_error
	}
	price = aes.GetDecryptedPrice(eprice)
	megaMediaLogger.Log(logger.INFO, "mega win price: ", price)

	OkRep(w)
	return

end_error:
	ErrorRep(w)
	return
}