package filter

import (
	"common"
	"sync"
)

type BasicFilter struct{}

// type RegionFilter struct{}
type MediaFilter struct{}
type TagFilter struct{}
type RetargetFilter struct{}

func NewBasicFilter() *BasicFilter {
	return &BasicFilter{}
}

func (f BasicFilter) Do(ad *common.Ad, req *common.BidRequest) (int, bool) {
	if !ad.Active {
		return 0, false
	}
	if len(req.Slots) == 0 {
		return 0, false
	}
	/* Now, we only deal with situation len(slots) == 1 */
	if req.Slots[0].W != ad.W || req.Slots[0].H != ad.H {
		return 0, false
	}
	return 0, true
}

func NewMediaFilter() *MediaFilter {
	return &MediaFilter{}
}

func (f MediaFilter) Do(ad *common.Ad, req *common.BidRequest) (int, bool) {
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
var mediaFilter MediaFilter
var tagFilter TagFilter
var retargetFilter RetargetFilter

func Init() {
	once.Do(func() {
		GFilterList = NewFilterList()
		GFilterList.Add(basicFilter)
		GFilterList.Add(mediaFilter)
		GFilterList.Add(tagFilter)
		GFilterList.Add(retargetFilter)
	})
}
