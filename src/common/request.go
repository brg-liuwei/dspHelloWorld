package common

import ()

type BidRequest struct {
	Id       string
	Mobile   bool
	Ip       string
	Did      string // SHA1(IMEI) for mobile
	Dpid     string // SHA1(AndroidID)
	Mac      string // SHA1(mac)
	Region   string
	Uuid     string
	Ua       string // user-agent
	Os       string
	Browser  string
	Url      string
	Domain   string
	Language string
	Media    string // "google,taobao,mango,..."
	Slots    []AdSlotType
	// esd, test, ping
	AdxId Adx

	User UserType
	Site SiteType
	App  AppType

	// TODO:
	// cookie string
}

type AdSlotType struct {
	ImpId      string
	BidFloor   int // 微分(0.0001 RMB/CPM)
	W          int
	H          int
	Visibility int
	CatOut     map[string]bool
	AttrOut    map[string]bool
	Instl      bool // 是否全插屏广告（移动端）
	Splash     bool // 是否开屏广告（移动端）
}

type GenderType int

const (
	GenderUnknown GenderType = iota
	GenderMale
	GenderFemale
)

type UserType struct {
	Gender GenderType
	Age    int // 0: unknown
	Cat    string
}

type SiteType struct {
	Cat     string
	Quality int
}

type AppType struct {
	Id      string
	Quality int
}
