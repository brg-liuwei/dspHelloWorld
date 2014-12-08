package filter

import (
	"common"
	"fmt"
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
			fmt.Println("ad", ad.Id, " ******* not pass filter list")
			return 0, false
		}
	}

	fmt.Println("ad", ad.Id, " pass filter list")

	/* last filter */
	//if basePrice, err := common.GOrderContainer.FindPrice(ad.OrderId, common.MANGO); err != nil {
	if basePrice, err := common.GOrderContainer.FindPrice(ad.OrderId, common.MEGAMEDIA); err != nil {
		fmt.Println("find basePrice err: ", err)
		return 0, false
	} else {
		/***********  fix libo's bug  ************************/
		if basePrice == 0 {
			fmt.Println("filter.go DoFilter, need to fix libo bug")
			basePrice = 7000
		}

		bidPrice := (basePrice + 1) * (100 + rate)
		bidPrice /= 100
		fmt.Println("basePrice: ", basePrice, " bidPrice: ", bidPrice, " floor: ", req.Slots[0].BidFloor)
		if len(req.Slots) != 0 && bidPrice >= req.Slots[0].BidFloor {
			return bidPrice, true
		} else {
			fmt.Println("filter.go: req.Slots len = ", len(req.Slots), " bid Price = ", bidPrice, " req.Slots[0].BidFloor: ", req.Slots[0].BidFloor)
		}
	}
	return 0, false
}
