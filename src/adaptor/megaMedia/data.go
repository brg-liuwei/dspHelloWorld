package megaMedia

import (
	"aes"
	"logger"

	"net/http"
	"strconv"
)

var megaMediaLogger *logger.Log

func Init(path string) {
	aes.Init()
	defer aes.Destroy()
	aes.SetKey("haha, I will never push this key to my github")

	megaMediaLogger = logger.NewLog(path)

	http.HandleFunc("/yesky/click", YeskyClickHandler)
	http.HandleFunc("/yesky/display", YeskyDisplayHandler)
	http.HandleFunc("/yesky/bid", YeskyBidHandler)
	http.HandleFunc("/yesky/win", YeskyWinHandler)
}

func ImpIdEncoding(id int32, aid int32) string {
	var code int64
	code = (int64(id) << 32) | (int64(aid) & 0xFFFFFFFF)
	return strconv.FormatInt(code, 10)
}

func ImpIdDecoding(imp string) (id int32, aid int32) {
	if code, err := strconv.ParseInt(imp, 10, 64); err == nil {
		id = int32((code >> 32) & 0xFFFFFFFF)
		aid = int32(code & 0xFFFFFFFF)
	}
	return
}
