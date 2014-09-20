package mango

import ()

type AdPosition int

const (
	TOPVIEW    AdPosition = 1 // 顶端可见
	BOTTOMVIEW                // 底端可见
	TOPROLL                   // 顶端随滚动条滚动可见
	BOTTOMROLL                // 底端随滚动条滚动可见
)

type AdType string

const (
	TEXT  AdType = "1"
	PIC          = "2"
	GIF          = "3"
	HTML         = "4"
	MRAID        = "5"
	VIDEO        = "6"
)

type AdAttr string

const (
	AutoAudio           AdAttr = "1" // Auto Play
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
	Other                   = "018"
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
	Other = 10
)
