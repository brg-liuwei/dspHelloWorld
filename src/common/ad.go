package common

import ()

type AdType int

const (
	AdtypeUnknown AdType = iota
	Banner
	Video
)

func (t *AdType) String() string {
	switch *t {
	case Banner:
		return "BANNER"
	case Video:
		return "VIDEO"
	default:
		return "UnknownAdType"
	}
}

type MimeType int

const (
	MimeUnknown MimeType = iota
	MimeVideo
	MimeX_FLV
)

func (t *MimeType) String() string {
	switch *t {
	case MimeVideo:
		return "Video"
	case MimeX_FLV:
		return "X-FLV"
	default:
		return "UnkownMimeType"
	}
}

type Attribute int

const (
	UNKOWN Attribute = iota
	TEXT
	PIC
	FLASH
	VIDEO
	TEXTLINK
	IFRAME
	JS
	HTML
	DYNAMICFLASH
)

type Ad struct {
	Id       string //ok
	Type     AdType // ok
	W        int    //ok
	H        int    //ok
	Priority int    /* default value: 0 */ // ok
	/* Attributes */
	Category        []int     /* category of ad */           //ok
	CategoryProduct []int     /* category of ad producer */  // ok
	Attr            Attribute /* attribute of ad creative */ // ok
	/* urls */
	UrlDisplayMonitor string
	UrlClickMonitor   string
	UrlThirdMonitor   string /* third party monitor url */ // ok
	UrlLanding        string /* landing page url */        //ok
	/* Templete */
	TmpName     string /* name of ad templete */                   // ok
	HtmlSnippet string /* html snippet code of dynamic creative */ //ok
	/* used for video */
	Mime        MimeType // ok
	Duration    int      /* time duration of video, 0 for banner */   //ok
	Ch1         string   /* first level channel of video creative */  // ok
	Ch2         string   /* second level channel of video creative */ // ok
	OrderId     string   // ok
	UrlCreative string   //ok
	CreativeId  int      // ok
	Active      bool     // ok

	// url溢价
	UrlIn    map[string]bool
	UrlOut   map[string]bool
	UrlPrice map[string]int

	// ad place溢价
	SlotIn    map[string]bool
	SlotOut   map[string]bool
	SlotPrice map[string]int
}

type Adx int

const (
	GOOGLE   Adx = iota + 1 // 1
	ALLYES                  // 2
	TAOBAO                  // 3
	SINA                    // 4
	TENCENT                 // 5
	SOHU                    // 6
	MIAOZHEN                // 7
	MANGO                   // 8
	NADX
)

type FreqCtrlType int

const (
	BYORDER FreqCtrlType = iota
	BYCAMPAIGN
)

type Order struct {
	Id         string
	CampaignId string
	AdOwnerId  string

	AdxList  []Adx /* Adx deliver list */
	MaxPrice [NADX]int

	/* ctrl data */
	Record      bool /* if record order-adx bid_num, succ_num and costs */
	MarkDynamic bool /* if mark dynamic creative order */

	/* frequency control */
	FreqCtrl        FreqCtrlType
	FreqCtrlTimeval int /* freq ctrl timeval, unit: seconds(86400 or 3600) */
	FreqCount       int

	LTime int // deliver time (24 bit mask)

	UseClassid  bool /* 使用动态素材标签定向 */
	UseRetarget bool /* 使用回头客定向 */

	AdvertX         int  /* 广告主暗扣比例系数（单位： %） */
	AdvertXUseMulti bool /* 是否使用乘法计算暗扣系数（默认为true） */
	SysX            int  /* 系统暗扣系数 */
	SysXUseMulti    bool /* 是否使用乘法计算系统暗扣系数（默认为true） */

	/* filters */
	AnonymousMedia int /* 匿名媒体定向设置, 0:不限；1:定向匿名媒体；2:排除匿名媒体 */
	Gender         int /* 性别定向: 0:all; 1:男; 2:女 */

	CountCost int /* 订单投放配额：微分 */
	Active    bool
}
