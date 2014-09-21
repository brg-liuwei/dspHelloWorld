package common

import ()

type BidAd struct {
	AdId           string
	OrderId        string
	ImpId          string
	Price          int
	WinUrl         string // win notice url
	Adm            string // html snippet or url
	W              int
	H              int
	DisplayMonitor string
	ClickMonitor   string // click monitor url
	LandingPage    string
	Domain         string
	Ext            string
}

type BidResponse struct {
	ReqId string
	BidId string
	Ads   []BidAd
}
