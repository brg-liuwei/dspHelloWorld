package mango

import (
	"net/http"

	"logger"
)

type AdPosition int

const (
	AdPosUnknown AdPosition = iota // 0
	TOPVIEW                        // 顶端可见
	BOTTOMVIEW                     // 底端可见
	TOPROLL                        // 顶端随滚动条滚动可见
	BOTTOMROLL                     // 底端随滚动条滚动可见
)

type AdType string

const (
	TEXT  AdType = "1"
	PIC          = "2"
	GIF          = "3"
	HTML         = "4"
	MRAID        = "5"
	VIDEO        = "6"
	FLASH        = "7"
)

type AdAttr string

const (
	AdAttrUnknown       AdAttr = "0"
	AutoAudio                  = "1" // Auto Play
	ClickAudio                 = "2"
	AutoExpandable             = "3"
	ClickExpandable            = "4"
	RollExpandable             = "5"
	AutoBannerVideo            = "6"
	ClickBannerVideo           = "7"
	Pop                        = "8"
	SuggestitiveImagery        = "9"  // provocative content
	Shaky                      = "10" // 闪烁
	Surveys                    = "11"
	TextOnly                   = "12"
	UserInteractive            = "13"
	WindowsDialogStyle         = "14"
	HasAudioOnOffButton        = "15"
	AdCanBeSkipped             = "16"
)

type AdCategory string

const (
	Game         AdCategory = "001" // 游戏类
	Clothes                 = "002" // 服饰类
	Cos                     = "003" // 日用化妆品类
	Net                     = "004" // 网络服务类
	Person                  = "005" // 个人用品类
	Retail                  = "006" // 零售类
	Amuse                   = "007" // 娱乐类
	Education               = "008" // 教育类
	Decorate                = "009" // 家居装饰类
	Food                    = "010" // 食品饮料类
	Traffic                 = "011" // 交通类
	IT                      = "012" // IT产品类
	Electronic              = "013" // 电子类
	Medical                 = "014" // 医疗类
	Finacial                = "015" // 金融类
	Operator                = "016" // 运营商类
	RealProperty            = "017" // 房地产类
	AdCatOther              = "018"
)

type Impression struct {
	ImpId       string
	BidFloor    int
	BidFloorCur string     // RMB = "CNY" 币种
	W           int        // Width
	H           int        // Height
	Pos         AdPosition // position of ad place
	Btype       []AdType   // deny ad type
	Battr       []AdAttr   // deny ad attr
	Instl       bool       // 是否全插屏ad
	Splash      bool       // 是否开屏ad
}

func (imp *Impression) Assign(m *map[string]interface{}) bool {
	if !imp.SetImpId(m) {
		return false
	}
	imp.SetBidFloor(m)
	imp.SetBidFloorCur(m)
	imp.SetW(m)
	imp.SetH(m)
	imp.SetPos(m)
	imp.SetBtype(m)
	imp.SetBattr(m)
	imp.SetInstl(m)
	imp.SetSplash(m)
	return true
}

func (imp *Impression) SetImpId(m *map[string]interface{}) bool {
	if v, ok := (*m)["impid"]; !ok {
		return false
	} else if imp.ImpId, ok = v.(string); !ok {
		return false
	} else {
		return true
	}
}

func (imp *Impression) SetBidFloor(m *map[string]interface{}) {
	if v, ok := (*m)["bidfloor"]; ok {
		f64, ok := v.(float64)
		if !ok {
			println("json bidfloor type error")
		} else {
			imp.BidFloor = int(f64)
		}
	}
}

func (imp *Impression) SetBidFloorCur(m *map[string]interface{}) {
	if v, ok := (*m)["bidfloorcur"]; ok {
		imp.BidFloorCur, _ = v.(string)
	}
}

func (imp *Impression) SetW(m *map[string]interface{}) {
	if v, ok := (*m)["w"]; ok {
		f64, ok := v.(float64)
		if !ok {
			println("json w type error")
		} else {
			imp.W = int(f64)
		}
	}
}

func (imp *Impression) SetH(m *map[string]interface{}) {
	if v, ok := (*m)["h"]; ok {
		f64, ok := v.(float64)
		if !ok {
			println("json h type error")
		} else {
			imp.H = int(f64)
		}
	}
}

func (imp *Impression) SetPos(m *map[string]interface{}) {
	if v, ok := (*m)["pos"]; ok {
		f64, ok := v.(float64)
		if !ok {
			println("json h type error")
		} else {
			imp.Pos = AdPosition(int(f64))
		}
	}
}

func (imp *Impression) SetBtype(m *map[string]interface{}) {
	if v, ok := (*m)["btype"]; ok {
		if array, ok := v.([]interface{}); ok {
			imp.Btype = make([]AdType, 0, len(array))
			for _, elem := range array {
				if e, ok := elem.(string); ok {
					imp.Btype = append(imp.Btype, AdType(e))
				}
			}
		}
	}
}

func (imp *Impression) SetBattr(m *map[string]interface{}) {
	if v, ok := (*m)["battr"]; ok {
		if array, ok := v.([]interface{}); ok {
			imp.Battr = make([]AdAttr, 0, len(array))
			for _, elem := range array {
				if e, ok := elem.(string); ok {
					imp.Battr = append(imp.Battr, AdAttr(e))
				}
			}
		}
	}
}

func (imp *Impression) SetInstl(m *map[string]interface{}) {
	if v, ok := (*m)["instl"]; ok {
		if instl, ok := v.(float64); ok {
			switch int(instl) {
			case 0:
				imp.Instl = false // 非插屏广告
			default:
				imp.Instl = true // 插屏广告
			}
		}
	}
}

func (imp *Impression) SetSplash(m *map[string]interface{}) {
	if v, ok := (*m)["splash"]; ok {
		if splash, ok := v.(float64); ok {
			switch int(splash) {
			case 0:
				imp.Splash = false // 非开屏广告
			case 1:
				imp.Splash = true
			}
		}
	}
}

type AppCat string

const (
	// prefix: A means Android, I means iOS
	Abook    AppCat = "60001"
	Abus            = "60002" // 商务
	Aani            = "60003" // 动漫
	Acom            = "60004" // 通讯
	Aedu            = "60005" // 教育
	Aamu            = "60006" // 娱乐
	Afin            = "60007" // 财务
	Agam            = "60008" // 游戏
	Ahea            = "60009" // 健康与健身
	Asof            = "60010" // 软件与演示
	Afas            = "60011" // 生活与时尚
	Awal            = "60012" // 动态壁纸
	Amedia          = "60013" // 媒体与视频
	Amedical        = "60014" // 医疗
	Amus            = "60015" // 音乐
	Anews           = "60016" // 新闻
	Aperson         = "60017" // 个性化
	Apho            = "60018" // 摄影
	Aeff            = "60019" // 效率
	Abuy            = "60020" // 购物
	Asoc            = "60021" // 社交
	Aspo            = "60022" // 体育
	Atool           = "60023" // 工具
	Atra            = "60024" // 交通
	Atour           = "60025" // 旅游
	Awea            = "60026" // 天气
	Apart           = "60027" // 小部件
	Agam1           = "60028" // 街机动作类
	Agam2           = "60029" // 解谜问答类
	Agam3           = "60030" // 扑克类
	Agam4           = "60031" // 休闲类
	Agam5           = "60032" // 动态壁纸类
	Agam6           = "60033" // 竞速类
	Agam7           = "60034" // 体育类
	Agam8           = "60035" // 小部件类
	// 60501 ~ 60570为ios类
	Ibook = "60051"
)

type Application struct {
	Aid         string
	Name        string
	Cat         []AppCat // category
	Ver         string   // version of app
	Bundle      string   // BundleID
	Itid        string   // iOS App iTunes ID
	Paid        bool     // is app paid of free ?
	Storeurl    string   // download url of app market
	Keywords    string   // keywords of app, separated by ","
	PublisherId string
	Publisher   string
}

func (app *Application) Assign(m *map[string]interface{}) bool {
	if !app.SetAid(m) {
		return false
	}
	app.SetName(m)
	app.SetCat(m)
	app.SetVer(m)
	app.SetBundle(m)
	app.SetItid(m)
	app.SetPaid(m)
	app.SetStoreurl(m)
	app.SetKeywords(m)
	app.SetPublisherId(m)
	app.SetPublisher(m)
	return true
}

func (app *Application) SetAid(m *map[string]interface{}) bool {
	if v, ok := (*m)["aid"]; !ok {
		return false
	} else if app.Aid, ok = v.(string); !ok {
		return false
	} else {
		return true
	}
}

func (app *Application) SetName(m *map[string]interface{}) {
	if v, ok := (*m)["name"]; ok {
		app.Name, _ = v.(string)
	}
}

func (app *Application) SetCat(m *map[string]interface{}) {
	if v, ok := (*m)["cat"]; ok {
		if array, ok := v.([]interface{}); ok {
			app.Cat = make([]AppCat, 0, len(array))
			for _, elem := range array {
				if e, ok := elem.(string); ok {
					app.Cat = append(app.Cat, AppCat(e))
				}
			}
		}
	}
}

func (app *Application) SetVer(m *map[string]interface{}) {
	if v, ok := (*m)["ver"]; ok {
		app.Ver, _ = v.(string)
	}
}

func (app *Application) SetBundle(m *map[string]interface{}) {
	if v, ok := (*m)["bundle"]; ok {
		app.Bundle, _ = v.(string)
	}
}

func (app *Application) SetItid(m *map[string]interface{}) {
	if v, ok := (*m)["itid"]; ok {
		app.Itid, _ = v.(string)
	}
}

func (app *Application) SetPaid(m *map[string]interface{}) {
	if v, ok := (*m)["paid"]; ok {
		paid, _ := v.(float64)
		switch int(paid) {
		case 0:
			app.Paid = false
		default:
			app.Paid = true
		}
	}
}

func (app *Application) SetStoreurl(m *map[string]interface{}) {
	if v, ok := (*m)["storeurl"]; ok {
		app.Storeurl, _ = v.(string)
	}
}

func (app *Application) SetKeywords(m *map[string]interface{}) {
	if v, ok := (*m)["Keywords"]; ok {
		app.Keywords, _ = v.(string)
	}
}

func (app *Application) SetPublisherId(m *map[string]interface{}) {
	if v, ok := (*m)["Pid"]; ok {
		app.PublisherId, _ = v.(string)
	}
}

func (app *Application) SetPublisher(m *map[string]interface{}) {
	if v, ok := (*m)["pub"]; ok {
		app.Publisher, _ = v.(string)
	}
}

type ConnType int

const (
	ConnUnkown ConnType = iota // 0
	ConnWifi
	ConnNG // 蜂窝数据网络未知几G
	Conn2G
	Conn3G
	Conn4G
)

type DeviceType int

const (
	DevUnkown     DeviceType = iota // 0
	DevIphone                       // iPhone
	DevAndroid                      // android phone
	DevIpad                         // iPad
	DevWphone                       // windows phone
	DevAndroidPad                   // android pad
	DevTv                           // 智能电视
)

type Device struct {
	Did         string     // SHAI(IMEI)
	Dpid        string     // Android Id or IDFA
	Mac         string     // SHA1(mac address)
	Ua          string     // User-Agent string of brower
	Ip          string     // ip
	Country     string     // ISO 3166-2
	Carrier     string     // 设备使用的运营商
	Language    string     // 语言
	Maker       string     // 设备制造商
	Model       string     // 设备型号
	Os          string     // 操作系统
	Osv         string     // Os version
	CType       ConnType   // 设备联网方式
	DType       DeviceType // 设备类型
	Loc         string     // 设备经纬度，逗号分隔，如38.04,114.50
	Sw          int        // 屏幕分辨率宽度像素数
	Sh          int        // 屏幕分辨率高度像素数
	Orientation int        // 屏幕方向： 1-竖向；2-横向
}

func (dev *Device) Assign(m *map[string]interface{}) {
	dev.SetDid(m)
	dev.SetDpid(m)
	dev.SetMac(m)
	dev.SetUa(m)
	dev.SetIp(m)
	dev.SetCountry(m)
	dev.SetCarrier(m)
	dev.SetLanguage(m)
	dev.SetMaker(m)
	dev.SetModel(m)
	dev.SetOs(m)
	dev.SetOsv(m)
	dev.SetCType(m)
	dev.SetDType(m)
	dev.SetLoc(m)
	dev.SetSw(m)
	dev.SetSh(m)
	dev.SetOrientation(m)
}

func (dev *Device) SetDid(m *map[string]interface{}) {
	if v, ok := (*m)["did"]; ok {
		dev.Did, _ = v.(string)
	}
}

func (dev *Device) SetDpid(m *map[string]interface{}) {
	if v, ok := (*m)["did"]; ok {
		dev.Did, _ = v.(string)
	}
}

func (dev *Device) SetMac(m *map[string]interface{}) {
	if v, ok := (*m)["mac"]; ok {
		dev.Mac, _ = v.(string)
	}
}

func (dev *Device) SetUa(m *map[string]interface{}) {
	if v, ok := (*m)["ua"]; ok {
		dev.Ua, _ = v.(string)
	}
}

func (dev *Device) SetIp(m *map[string]interface{}) {
	if v, ok := (*m)["ip"]; ok {
		dev.Ip, _ = v.(string)
	}
}

func (dev *Device) SetCountry(m *map[string]interface{}) {
	if v, ok := (*m)["country"]; ok {
		dev.Country, _ = v.(string)
	}
}

func (dev *Device) SetCarrier(m *map[string]interface{}) {
	if v, ok := (*m)["carrier"]; ok {
		dev.Carrier, _ = v.(string)
	}
}

func (dev *Device) SetLanguage(m *map[string]interface{}) {
	if v, ok := (*m)["language"]; ok {
		dev.Language, _ = v.(string)
	}
}

func (dev *Device) SetMaker(m *map[string]interface{}) {
	if v, ok := (*m)["make"]; ok {
		dev.Maker, _ = v.(string)
	}
}

func (dev *Device) SetModel(m *map[string]interface{}) {
	if v, ok := (*m)["model"]; ok {
		dev.Model, _ = v.(string)
	}
}

func (dev *Device) SetOs(m *map[string]interface{}) {
	if v, ok := (*m)["os"]; ok {
		dev.Os, _ = v.(string)
	}
}

func (dev *Device) SetOsv(m *map[string]interface{}) {
	if v, ok := (*m)["osv"]; ok {
		dev.Osv, _ = v.(string)
	}
}

func (dev *Device) SetCType(m *map[string]interface{}) {
	if v, ok := (*m)["connectiontype"]; ok {
		if conn, ok := v.(float64); ok {
			dev.CType = ConnType(int(conn))
		}
	}
}

func (dev *Device) SetDType(m *map[string]interface{}) {
	if v, ok := (*m)["devicetype"]; ok {
		if dtype, ok := v.(float64); ok {
			dev.DType = DeviceType(int(dtype))
		}
	}
}

func (dev *Device) SetLoc(m *map[string]interface{}) {
	if v, ok := (*m)["loc"]; ok {
		dev.Loc, _ = v.(string)
	}
}

func (dev *Device) SetSw(m *map[string]interface{}) {
	if v, ok := (*m)["sw"]; ok {
		if sw, ok := v.(float64); ok {
			dev.Sw = int(sw)
		}
	}
}

func (dev *Device) SetSh(m *map[string]interface{}) {
	if v, ok := (*m)["sh"]; ok {
		if sh, ok := v.(float64); ok {
			dev.Sh = int(sh)
		}
	}
}

func (dev *Device) SetOrientation(m *map[string]interface{}) {
	if v, ok := (*m)["orientation"]; ok {
		if orie, ok := v.(float64); ok {
			dev.Orientation = int(orie)
		}
	}
}

type ClickType int

const (
	CTypeUnkown ClickType = iota // 0
	SendMail
	DownloadApp
	OpenWeb
	SendMsg
	CallPhone
	GoToAppStore
	OpenMap
	CliTypeOther = 10
)

var mangoLogger *logger.Log

func Init(path string) {
	mangoLogger = logger.NewLog(path)

	http.HandleFunc("/mango/bid", MangoBidHandler)
	http.HandleFunc("/mango/win", MangoWinHandler)
	http.HandleFunc("/mango/click", MangoClickHandler)
	http.HandleFunc("/mango/display", MangoDisplayHandler)
}
