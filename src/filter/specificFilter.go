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
		return 0, false
	}
	if req.Slots == nil || len(req.Slots) == 0 {
		return 0, false
	}
	/* Now, we only deal with situation len(slots) == 1 */
	if req.Slots[0].W != ad.W || req.Slots[0].H != ad.H {
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
			return 0, false
		}
	}
	if len(ad.UrlOut) != 0 {
		if _, ok := ad.UrlOut[url]; ok {
			return 0, false
		}
	}
	rate, _ := ad.UrlPrice[url]
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
