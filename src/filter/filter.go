package filter

import (
	"common"
)

type Filter interface {
	Do(ad *common.Ad, req *common.BidRequest) (rate int, canDeliver bool)
}

type FilterList struct {
	Filters []Filter
}

var GFilterList *FilterList

func NewFilterList() *FilterList {
	return &FilterList{
		Filters: make([]Filter, 0, 8),
	}
}

func (fl *FilterList) Add(f Filter) {
	fl.Filters = append(fl.Filters, f)
}

/*
int: bidPrice
bool: can deliver?
*/
func (fl *FilterList) DoFilter(ad *common.Ad, req *common.BidRequest) (int, bool) {
	rate := 0
	for _, f := range fl.Filters {
		if r, ok := f.Do(ad, req); ok {
			rate += r
		} else {
			return 0, false
		}
	}

	/* last filter */
	//if basePrice, err := common.GOrderContainer.FindPrice(ad.OrderId, common.MANGO); err != nil {
	if basePrice, err := common.GOrderContainer.FindPrice(ad.OrderId, common.MEGAMEDIA); err != nil {
		return 0, false
	} else {
		bidPrice := basePrice * (100 + rate)
		bidPrice /= 100
		if len(req.Slots) != 0 && bidPrice >= req.Slots[0].BidFloor {
			return bidPrice, true
		}
	}
	return 0, false
}
