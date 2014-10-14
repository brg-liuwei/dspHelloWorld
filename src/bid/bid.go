package bid

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"strconv"

	"common"
	"filter"
)

func Bid(req *common.BidRequest) *common.BidResponse {
	ids := make([]string, 0, 4)
	prices := make([]int, 0, 4)
	fmt.Println("\nin bid len(ads): ", len(common.GAdContainer.Ads))
	func() {
		c := common.GAdContainer
		c.Lock.RLock()
		defer c.Lock.RUnlock()
		idx := rand.Int()
		for i := 0; i < len(c.Ads) && i < 256; i++ {
			idx = idx % len(c.Ads)
			if price, ok := filter.GFilterList.DoFilter(&c.Ads[idx], req); ok {
				ids = append(ids, c.Ads[idx].Id)
				prices = append(prices, price)
				if len(ids) >= 3 {
					return
				}
			}
		}
	}()
	if len(ids) == 0 {
		return GenEmptyBidResponse(req)
	} else {
		idx := rand.Int() % len(ids)
		fmt.Println("idx = ", idx, "ids[idx] = ", ids[idx])
		ad, err := common.GAdContainer.Find(ids[idx])
		if err != nil {
			return GenEmptyBidResponse(req)
		}
		return GenBidResponse(req, &ad, prices[idx])
	}
}

func genMd5(str string) string {
	md5hash := md5.New()
	io.WriteString(md5hash, str)
	return fmt.Sprintf("%x", md5hash.Sum(nil))
}

func genBidId() string {
	return genMd5(strconv.FormatInt(rand.Int63(), 10))
}

func GenEmptyBidResponse(req *common.BidRequest) *common.BidResponse {
	return &common.BidResponse{
		ReqId: req.Id,
		BidId: genBidId(),
		Ads:   nil,
	}
}

func GenBidResponse(req *common.BidRequest, ad *common.Ad, price int) *common.BidResponse {
	if len(req.Slots) == 0 {
		return GenEmptyBidResponse(req)
	}
	rep := new(common.BidResponse)
	rep.ReqId = req.Id
	rep.BidId = genBidId()
	rep.Ads = make([]common.BidAd, 1)

	rep.Ads[0].AdId = ad.Id
	rep.Ads[0].OrderId = ad.OrderId
	rep.Ads[0].ImpId = req.Slots[0].ImpId
	rep.Ads[0].Price = price
	rep.Ads[0].WinUrl = common.WinUrl
	//rep.Ads[0].Adm = ad.HtmlSnippet
	rep.Ads[0].Adm = ad.UrlCreative
	rep.Ads[0].W = ad.W
	rep.Ads[0].H = ad.H
	rep.Ads[0].DisplayMonitor = ad.UrlDisplayMonitor
	rep.Ads[0].ClickMonitor = ad.UrlClickMonitor
	rep.Ads[0].LandingPage = ad.UrlLanding
	//rep.Ads[0].Domain = //...
	//rep.Ads[0].Ext.xxx = //...
	return rep
}
