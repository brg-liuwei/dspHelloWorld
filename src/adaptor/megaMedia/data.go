package megaMedia

import (
	"logger"
	"strconv"
)

var megaMediaLogger *logger.Log

func Init(path string) {
	megaMediaLogger = logger.NewLog(path)
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
