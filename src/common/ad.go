package common

import ()

type AdType int

const (
	Banner AdType = iota
	Video
)

func (t *AdType) String() string {
	switch t {
	case Banner:
		return "BANNER"
	case Video:
		return "VIDEO"
	default:
		return "UnknownAdType"
	}
}

type AdSize struct {
	X int
	Y int
}

type MimeType int

const (
	X_FLV MimeType = iota
)

func (t *MimeType) String() string {
	switch t {
	case X_FLV:
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
	Id       int
	Type     AdType
	Size     AdSize
	Priority int /* default value: 0 */

	/* Attributes */
	Category        []int     /* category of ad */
	CategoryProduct []int     /* category of ad producer */
	Attr            Attribute /* attribute of ad creative */

	/* urls */
	UrlThird   string /* third party monitor url */
	UrlLanding string /* landing page url */

	/* Templete */
	TmpName     string /* name of ad templete */
	HtmlSnippet string /* html snippet code of dynamic creative */

	/* used for video */
	Mime     MimeType
	Duration int /* time duration of video, 0 for banner */
	Ch1          /* first level channel of video creative */
	Ch2          /* second level channel of video creative */

	OrderId    int
	CreativeId int
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
	MANGO
	NADX
)

type FreqCtrlType int

const (
	BYORDER FreqCtrlType = iota
	BYCAMPAIGN
)

type Order struct {
	Id         int
	CampaignId int
	AdOwnerId  int

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
}
