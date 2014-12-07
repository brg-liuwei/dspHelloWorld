package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"common"
	level "github.com/brg-liuwei/golevel"
	redis "github.com/gosexy/redis"
	"logger"
)

var cmdTable string
var ErrNoCmd error = errors.New("No Cmd")

type Commander struct {
	host string
	port uint
	rkey string // key of redis queue
	cli  *redis.Client
	err  error
}

func NewCommander(host string, port uint, rkey string) *Commander {
	c := redis.New()
	err := c.Connect(host, port)
	return &Commander{host: host,
		port: port,
		rkey: rkey,
		cli:  c,
		err:  err}
}

func (c *Commander) Repair() bool {
	if c.err == nil {
		return true
	}
	if c.err = c.cli.Connect(c.host, c.port); c.err != nil {
		return false
	}
	return true
}

func (c *Commander) GetCmd() (cmd string, err error) {
	if !c.Repair() {
		return "", c.err
	}
	ss, e := c.cli.BLPop(600, c.rkey)
	switch e {
	case nil:
		if len(ss) != 2 {
			managerLogger.Log(logger.ERROR, "command len error: ", len(ss))
		} else {
			cmd = ss[1]
		}
	case redis.ErrNilReply:
		err = ErrNoCmd
	default:
		err = e
		c.Repair()
	}
	return
}

func CommanderRoutine(host string, port uint, rkey string) {
	level.Init(4)
	defer level.Cleanup()
	level.Open(&cmdTable)
	defer level.Close(&cmdTable)

	command := NewCommand()

	cNum := 100000000 // to keep cmd by order
	itr, err := level.NewIterator(&cmdTable)
	if err != nil {
		panic(err)
	}
	for itr.SeekToFirst(); itr.Valid(); itr.Next() {
		cmd := itr.Value()
		if command.Parse(cmd) {
			managerLogger.Log(logger.DEBUG, cmd)
			command.Execute()
		}
		cNum++
	}
	itr.Destroy()

	c := NewCommander(host, port, rkey)
	for {
		cmd, err := c.GetCmd()
		switch err {
		case ErrNoCmd:
			time.Sleep(5 * time.Second)
		case nil:
			if command.Parse(cmd) {
				// if execute cmd successful, we store this cmd
				if command.Execute() {
					key := strconv.Itoa(cNum)
					level.Put(&cmdTable, &key, &cmd)
					cNum++
				}
			}
		default:
			managerLogger.Log(logger.ERROR, "redis connect err:", err)
			c.Repair()
			time.Sleep(1 * time.Second)
		}
	}
}

type CmdType int

const (
	NilCmd   CmdType = iota // 0
	AddOrder                // 1
	AddAd                   // 2
	ModOrder                // 3
	ModAd                   // 4
	DelOrder                // 5
	DelAd                   // 6
)

type Command struct {
	Ctype    CmdType
	Cversion string
	Data     interface{}
}

func NewCommand() *Command {
	return &Command{Ctype: NilCmd,
		Cversion: "NoVersion"}
}

func (c *Command) Parse(jsonCmd string) bool {
	type Fmt struct {
		Oper_type string
		Fmt_ver   string
		Data      interface{}
	}
	var cmd Fmt
	dec := json.NewDecoder(strings.NewReader(jsonCmd))
	err := dec.Decode(&cmd)
	if err != nil {
		managerLogger.Log(logger.ERROR, "decode cmd err: ", jsonCmd, err)
		return false
	}
	if slice, ok := cmd.Data.([]interface{}); ok {
		if len(slice) != 1 {
			managerLogger.Log(logger.ERROR, "cmd Data slice len should be 1 ", jsonCmd)
			return false
		}
	}
	switch cmd.Oper_type {
	case "1":
		c.Ctype = AddOrder
	case "2":
		c.Ctype = AddAd
	case "3":
		c.Ctype = ModOrder
	// case "4":
	// 	c.Ctype = ModAd
	case "5":
		c.Ctype = DelOrder
	case "6":
		c.Ctype = DelAd
	default:
		managerLogger.Log(logger.ERROR, "cmd type error: ", jsonCmd, cmd.Oper_type)
		return false
	}
	c.Cversion = cmd.Fmt_ver
	c.Data = cmd.Data
	return true
}

func (c *Command) Execute() bool {
	switch c.Ctype {
	case AddOrder:
		return c.AddOrder()
	case AddAd:
		return c.AddAd()
	case ModOrder:
		return c.ModOrder()
	case DelOrder:
		return c.DelOrder()
	case DelAd:
		return c.DelAd()
	}
	managerLogger.Log(logger.ERROR, "Execute cmd type error: ", c.Ctype)
	return false
}

func (c *Command) ModOrder() bool {
	var data []interface{}
	if dataSlice, ok := c.Data.([]interface{}); ok {
		data = dataSlice
	} else if d, ok := c.Data.(map[string]interface{}); ok {
		data = make([]interface{}, 0, 1)
		data = append(data, d)
	} else {
		managerLogger.Log(logger.ERROR, "in ModOrder data err")
		return false
	}

	type OrderModFmt struct {
		Order_id   string
		Balance    string /* Balance是该订单还剩余的钱数 */
		Count_cost string /* Count_cost是订单总价 */
	}

	for _, d := range data {
		var orderFmt OrderModFmt
		b, e := json.Marshal(d)
		if e != nil {
			managerLogger.Log(logger.ERROR, "Marshal json err: ", e)
			continue
		}
		dec := json.NewDecoder(bytes.NewReader(b))
		if e = dec.Decode(&orderFmt); e != nil {
			managerLogger.Log(logger.ERROR, "mod order json decode err: ", e)
			continue
		}

		if cost, err := strconv.ParseFloat(orderFmt.Balance, 64); err == nil {
			common.GOrderContainer.SetCost(orderFmt.Order_id, int(cost))
		} else {
			managerLogger.Log(logger.ERROR, "in ModOrder Balance err:", err)
		}
	}
	// always return false, to tell manager do NOT save this cmd
	return false
}

func (c *Command) AddOrder() bool {
	var data []interface{}
	if dataSlice, ok := c.Data.([]interface{}); ok {
		data = dataSlice
	} else if d, ok := c.Data.(map[string]interface{}); ok {
		data = make([]interface{}, 0, 1)
		data = append(data, d)
	} else {
		managerLogger.Log(logger.ERROR, "in AddOrder data err")
		return false
	}

	type OrderAddFmt struct {
		Order_id      string
		Campaign_id   string
		Advertiser_id string
		Exchange      []string

		Max_price []string

		Record_data string
		Is_dm       string

		/* freq ctrl暂时不用 */

		Classid_use      string
		Retarget_use     string
		Advert_X         string
		Ad_X_is_Multi    string
		Sys_X            string
		Sys_is_Multi     string
		Anonymous_filter string
		Gender_filter    string
		Count_cost       string

		Active string
	}

	// len(data) should be 1
	for _, d := range data {
		var orderFmt OrderAddFmt
		b, e := json.Marshal(d)
		if e != nil {
			managerLogger.Log(logger.ERROR, "add order Marshal json err: ", e)
			return false
		}
		dec := json.NewDecoder(bytes.NewReader(b))
		if e := dec.Decode(&orderFmt); e != nil {
			managerLogger.Log(logger.ERROR, "add order json decode err: ", e)
			return false
		}

		var order common.Order
		func(dst *common.Order, src *OrderAddFmt) {
			dst.Id = src.Order_id
			dst.CampaignId = src.Campaign_id
			dst.AdOwnerId = src.Advertiser_id
			dst.AdxList = make([]common.Adx, 0, 8)
			for _, e := range src.Exchange {
				if i, err := strconv.Atoi(e); err == nil {
					dst.AdxList = append(dst.AdxList, common.Adx(i))
				} else {
					managerLogger.Log(logger.ERROR, "Exchange Code Error: ", e)
				}
			}

			sentinel := int(common.NADX) - 1
			if sentinel > len(src.Max_price) {
				sentinel = len(src.Max_price)
			}
			for i := 0; i < sentinel; i++ {
				// Adx编号是从1开始的
				dst.MaxPrice[i+1], _ = strconv.Atoi(src.Max_price[i])
			}

			if src.Record_data == "1" {
				dst.Record = true
			} else {
				dst.Record = false
			}
			if src.Is_dm == "1" {
				dst.MarkDynamic = true
			} else {
				dst.MarkDynamic = false
			}
			if src.Classid_use == "1" {
				dst.UseClassid = true
			} else {
				dst.UseClassid = false
			}
			if src.Retarget_use == "1" {
				dst.UseRetarget = true
			} else {
				dst.UseRetarget = false
			}
			dst.AdvertX, _ = strconv.Atoi(src.Advert_X)
			if src.Ad_X_is_Multi == "0" {
				dst.AdvertXUseMulti = false
			} else {
				dst.AdvertXUseMulti = true
			}
			dst.SysX, _ = strconv.Atoi(src.Sys_X)
			if src.Sys_is_Multi == "0" {
				dst.SysXUseMulti = false
			} else {
				dst.SysXUseMulti = true
			}
			dst.AnonymousMedia, _ = strconv.Atoi(src.Anonymous_filter)
			dst.Gender, _ = strconv.Atoi(src.Gender_filter)
			dst.CountCost, _ = strconv.Atoi(src.Count_cost)
			if src.Active == "1" {
				dst.Active = true
			} else {
				dst.Active = false
			}
		}(&order, &orderFmt)
		norder := common.GOrderContainer.Add(&order)
		managerLogger.Log(logger.INFO, "add order successfully: ", d, "norders: ", norder)
		return true
	}
	return false
}

func (c *Command) DelOrder() bool {
	var data []interface{}
	if dataSlice, ok := c.Data.([]interface{}); ok {
		data = dataSlice
	} else if d, ok := c.Data.(map[string]interface{}); ok {
		data = make([]interface{}, 0, 1)
		data = append(data, d)
	} else {
		managerLogger.Log(logger.ERROR, "in DelOrder data err")
		return false
	}

	// len(data) == 1
	for _, d := range data {
		if id, ok := d.(string); ok {
			norder := common.GOrderContainer.Del(id)
			managerLogger.Log(logger.INFO, "del order successfully: ", d, "orders: ", norder)
			return true
		} else {
			managerLogger.Log(logger.ERROR, "in DelOrder, data array fmt err")
			return false
		}
	}
	return false
}

func (c *Command) AddAd() bool {
	var data []interface{}
	if dataSlice, ok := c.Data.([]interface{}); ok {
		data = dataSlice
	} else if d, ok := c.Data.(map[string]interface{}); ok {
		data = make([]interface{}, 0, 1)
		data = append(data, d)
	} else {
		managerLogger.Log(logger.ERROR, "in AddAd data err")
		return false
	}

	type AdAddFmt struct {
		Ad_id              string
		Order_id           string
		Ad_type            string
		Duration           string
		Mimes              string
		Channel            string
		Cs                 string
		Size               string
		Tmp_name           string
		Priority           string
		Clickmonitor_url   string
		Displaymonitor_url string
		Thirdparty_url     string
		Landingpage_url    string
		Creative_url       string
		Html_snippet       string
		Exchange           []string // 在哪儿？
		Category           []string
		Category_product   []string
		Attribute          string
		Buyer_creative_id  string
		Advertiser_name    string
		Active             string

		Url_in    []string
		Url_out   []string
		Url_price []map[string]string

		Slotid_in    []string
		Slotid_out   []string
		Slotid_price []map[string]string
	}

	for _, d := range data {
		var adFmt AdAddFmt
		b, e := json.Marshal(d)
		if e != nil {
			managerLogger.Log(logger.ERROR, "Marshal json err: ", e)
			continue
		}
		dec := json.NewDecoder(bytes.NewReader(b))
		if e := dec.Decode(&adFmt); e != nil {
			continue
		}

		var ad common.Ad
		func(dst *common.Ad, src *AdAddFmt) {
			dst.Id = src.Ad_id
			dst.OrderId = src.Order_id
			switch src.Ad_type {
			case "1":
				dst.Type = common.AdType(common.Banner)
			case "2":
				dst.Type = common.AdType(common.Video)
			default:
				managerLogger.Log(logger.ERROR, "Ad_type unknow: ", src.Ad_type)
				dst.Type = common.AdType(common.AdtypeUnknown)
			}
			dst.Duration, _ = strconv.Atoi(src.Duration)
			switch src.Mimes {
			case "video":
				dst.Mime = common.MimeType(common.MimeVideo)
			case "x-flv":
				dst.Mime = common.MimeType(common.MimeX_FLV)
			default:
				dst.Mime = common.MimeType(common.MimeUnknown)
			}
			dst.Ch1 = src.Channel
			dst.Ch2 = src.Cs
			s := strings.SplitN(src.Size, "x", 2)
			if len(s) != 2 {
				s = strings.SplitN(src.Size, "X", 2)
			}
			if len(s) == 2 {
				w, _ := strconv.Atoi(s[0])
				h, _ := strconv.Atoi(s[1])
				if w > 0 && h > 0 {
					dst.W = w
					dst.H = h
				}
			}
			dst.TmpName = src.Tmp_name
			dst.Priority, _ = strconv.Atoi(src.Priority)
			dst.UrlDisplayMonitor = src.Displaymonitor_url
			dst.UrlClickMonitor = src.Clickmonitor_url
			dst.UrlThirdMonitor = src.Thirdparty_url
			dst.UrlLanding = src.Landingpage_url
			dst.UrlCreative = src.Creative_url
			dst.HtmlSnippet = src.Html_snippet
			dst.Category = make([]int, 0, len(src.Category))
			for _, cat := range src.Category {
				if c, err := strconv.Atoi(cat); err == nil {
					dst.Category = append(dst.Category, c)
				}
			}
			dst.CategoryProduct = make([]int, 0, len(src.Category_product))
			for _, cp := range src.Category_product {
				if c, err := strconv.Atoi(cp); err == nil {
					dst.CategoryProduct = append(dst.CategoryProduct, c)
				}
			}
			if attrNum, err := strconv.Atoi(src.Attribute); err != nil {
				dst.Attr = common.Attribute(attrNum)
			}
			dst.CreativeId, _ = strconv.Atoi(src.Buyer_creative_id)
			switch src.Active {
			case "0":
				dst.Active = false
			default:
				dst.Active = true
			}

			/* for url filter */
			if src.Url_out == nil {
				dst.UrlOut = nil
			} else {
				dst.UrlOut = make(map[string]bool)
				for _, url := range src.Url_out {
					dst.UrlOut[url] = true
				}
			}
			if src.Url_in == nil {
				dst.UrlIn = nil
				dst.UrlPrice = nil
			} else {
				dst.UrlIn = make(map[string]bool)
				for _, url := range src.Url_in {
					dst.UrlIn[url] = true
				}
				if src.Url_price != nil {
					dst.UrlPrice = make(map[string]int)
					for _, m := range src.Url_price {
						url := m["Url"]
						price, _ := strconv.Atoi(m["Price"])
						dst.UrlPrice[url] = price
					}
				}
			}

			/* for slot filter */
			if src.Slotid_out == nil {
				dst.SlotOut = nil
			} else {
				dst.SlotOut = make(map[string]bool)
				for _, slot := range src.Slotid_out {
					dst.SlotOut[slot] = true
				}
			}
			if src.Slotid_in == nil {
				dst.SlotIn = nil
				dst.SlotPrice = nil
			} else {
				dst.SlotIn = make(map[string]bool)
				for _, slot := range src.Slotid_in {
					dst.SlotIn[slot] = true
				}
				if src.Slotid_price != nil {
					dst.SlotPrice = make(map[string]int)
					for _, m := range src.Slotid_price {
						slot := m["Slotid"]
						price, _ := strconv.Atoi(m["Price"])
						dst.SlotPrice[slot] = price
					}
				}
			}
		}(&ad, &adFmt)
		if _, err := common.GOrderContainer.Find(ad.OrderId); err == nil {
			nad := common.GAdContainer.Add(&ad)
			managerLogger.Log(logger.INFO, "add ad successfully: ", d, ", ads: ", nad)
			return true
		} else {
			managerLogger.Log(logger.ERROR, "add ad error: ", d, ", err: ", err)
			return false
		}
	}
	return false
}

func (c *Command) DelAd() bool {
	var data []interface{}
	if dataSlice, ok := c.Data.([]interface{}); ok {
		data = dataSlice
	} else if d, ok := c.Data.(map[string]interface{}); ok {
		data = make([]interface{}, 0, 1)
		data = append(data, d)
	} else {
		managerLogger.Log(logger.ERROR, "in DelAd data err")
		return false
	}

	// len(data) == 1
	for _, d := range data {
		if id, ok := d.(string); ok {
			nad := common.GAdContainer.Del(id)
			managerLogger.Log(logger.INFO, "del ad successfully: ", d, "ads: ", nad)
			return true
		} else {
			managerLogger.Log(logger.ERROR, "in DelAd, data array fmt err")
			return false
		}
	}
	return false
}

var managerLogger *logger.Log

func Init(path string) {
	managerLogger = logger.NewLog(path)
	cmdTable = "/opt/dspHelloWorld/dspCmd"
}
