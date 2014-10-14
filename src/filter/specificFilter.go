package filter

import (
	"common"
	"sync"
)

type BasicFilter struct{}
type UrlFilter struct{}
type SlotFilter struct{}
type TagFilter struct{}
type RetargetFilter struct{}

func NewBasicFilter() *BasicFilter {
	return &BasicFilter{}
}

func (f BasicFilter) Do(ad *common.Ad, req *common.BidRequest) (int, bool) {
	if !ad.Active {
		println("basic filter ad active false: ", ad.Id)
		return 0, false
	}
	if req.Slots == nil || len(req.Slots) == 0 {
		println("basic filter: req slots 0")
		return 0, false
	}
	/* Now, we only deal with situation len(slots) == 1 */
	if req.Slots[0].W != ad.W || req.Slots[0].H != ad.H {
		println("basic filter: ad w,h error", ad.Id)
		return 0, false
	}
	return 0, true
}

func NewUrlFilter() *UrlFilter {
	return &UrlFilter{}
}

func (f UrlFilter) Do(ad *common.Ad, req *common.BidRequest) (int, bool) {
	url := req.Url
	if len(ad.UrlIn) != 0 {
		if _, ok := ad.UrlIn[url]; !ok {
			println("url filter: url in !ok:", ad.Id)
			return 0, false
		}
	}
	if len(ad.UrlOut) != 0 {
		if _, ok := ad.UrlOut[url]; ok {
			println("url filter: url out :", ad.Id)
			return 0, false
		}
	}
	rate, _ := ad.UrlPrice[url]
	println(ad.Id, " pass url filter, rate: ", rate)
	return rate, true
}

func NewSlotFilter() *SlotFilter {
	return &SlotFilter{}
}

func (f SlotFilter) Do(ad *common.Ad, req *common.BidRequest) (int, bool) {
	return 0, true
}

func NewTagFilter() *TagFilter {
	return &TagFilter{}
}

func (f TagFilter) Do(ad *common.Ad, req *common.BidRequest) (int, bool) {
	return 0, true
}

func NewRetargetFilter() *RetargetFilter {
	return &RetargetFilter{}
}

func (f RetargetFilter) Do(ad *common.Ad, req *common.BidRequest) (int, bool) {
	return 0, true
}

var once sync.Once

var basicFilter BasicFilter
var urlFilter UrlFilter
var slotFilter SlotFilter
var tagFilter TagFilter
var retargetFilter RetargetFilter

func Init() {
	once.Do(func() {
		GFilterList = NewFilterList()
		GFilterList.Add(basicFilter)
		GFilterList.Add(urlFilter)
		GFilterList.Add(slotFilter)
		GFilterList.Add(tagFilter)
		GFilterList.Add(retargetFilter)
	})
}
