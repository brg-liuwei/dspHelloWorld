package megaMedia

import (
	"fmt"
	"net/http"

	proto "code.google.com/p/goprotobuf/proto"

	"common"
	"logger"
)

func (resp *MgxBidResponse) ParseFromCommon(cb *common.BidResponse) {
	resp.Version = proto.Int32(1)
	resp.BidGuid = proto.String(cb.ReqId)
	if len(cb.Ads) != 0 {
		resp.Adslot = make([]*MgxBidResponse_AdSlot, len(cb.Ads))
		for i, s := range cb.Ads {
			id, aid := ImpIdDecoding(s.ImpId)
			resp.Adslot[i].Id = proto.Int32(id)
			resp.Adslot[i].MegaxAid = proto.Int32(aid)
			// s.price的单位是微分/CPM，这里需要转换成分/CPM
			resp.Adslot[i].MaxCpmMicros = proto.Int64(int64(s.Price / 1000))
			// TODO: add our buyer id, advertise id
			// resp.Adslot[i].BuyerId = proto.String(ourid)
			// resp.Adslot[i].AdvertiserId = proto.String(aderId)
			// resp.CreativeFormat = proto.Int32(): 1,图片， 2，flash，3，flv，4,html或者iframe代码

			//---------------- 注意 -------------
			//这个Adm代码需要实现点击宏
			//-----------------------------------
			resp.Adslot[i].CreativeContent = proto.String(s.Adm)

			resp.Adslot[i].CreativeWidth = proto.Int32(int32(s.W))
			resp.Adslot[i].CreativeHeight = proto.Int32(int32(s.H))
			resp.Adslot[i].ClickThroughUrl = proto.String(s.LandingPage)

			resp.Adslot[i].ClickTracking = make([]string, 1)
			resp.Adslot[i].ClickTracking[0] = s.ClickMonitor

			resp.Adslot[i].ImpressionTracking = make([]string, 1)
			resp.Adslot[i].ImpressionTracking[0] = s.DisplayMonitor

			resp.Adslot[i].ExtendData = proto.String(s.Ext)
		}
	}
}

func (resp *MgxBidResponse) Response(w http.ResponseWriter) {
	data, err := proto.Marshal(resp)
	if err != nil {
		megaMediaLogger.Log(logger.ERROR, "megaMedia proto marshal err: ", err)
		goto end
	} else {
		// TODO: log type
		megaMediaLogger.Log(logger.INFO, "bid log")
	}

end:
	w.Header().Set("Content-type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Write(data)
}
