package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"common"
	redis "github.com/gosexy/redis"
	"logger"
)

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
		cmd = ss[0]
	case redis.ErrNilReply:
		err = ErrNoCmd
	default:
		err = e
		c.Repair()
	}
	return
}

func CommanderRoutine(host string, port uint, rkey string) {
	c := NewCommander(host, port, rkey)
	command := NewCommand()
	for {
		cmd, err := c.GetCmd()
		switch err {
		case ErrNoCmd:
			time.Sleep(5 * time.Second)
		case nil:
			if command.Parse(cmd) {
				command.Execute()
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
	switch cmd.Oper_type {
	case "1":
		c.Ctype = AddOrder
	case "2":
		c.Ctype = AddAd
	// case "3":
	// 	c.Ctype = ModOrder
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

func (c *Command) Execute() {
	switch c.Ctype {
	case AddOrder:
		c.AddOrder()
	case AddAd:
		c.AddAd()
	case DelOrder:
		c.DelOrder()
	case DelAd:
		c.DelAd()
	}
}

func (c *Command) AddOrder() {
	data, ok := c.Data.([]interface{})
	if !ok {
		managerLogger.Log(logger.ERROR, "in AddOrder data err")
		return
	}
	type OrderAddFmt struct {
		Order_id      string
		Campaign_id   string
		Advertiser_id string
		Exchange      []string

		/* waiting to add more exchange */
		Max_price_mango string

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

	for _, d := range data {
		var orderFmt OrderAddFmt
		b, e := json.Marshal(d)
		if e != nil {
			managerLogger.Log(logger.ERROR, "Marshal json err: ", e)
			continue
		}
		dec := json.NewDecoder(bytes.NewReader(b))
		if e := dec.Decode(&orderFmt); e != nil {
			continue
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
			dst.MaxPrice[common.MANGO], _ = strconv.Atoi(src.Max_price_mango)
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
		common.GOrderContainer.Add(&order)
	}
}

func (c *Command) DelOrder() {
	data, ok := c.Data.([]interface{})
	if !ok {
		managerLogger.Log(logger.ERROR, "in DelOrder data err")
		return
	}
	for _, d := range data {
		if id, ok := d.(string); ok {
			common.GOrderContainer.Del(id)
		} else {
			managerLogger.Log(logger.ERROR, "in DelOrder, data array fmt err")
		}
	}
}

func (c *Command) AddAd() {
	data, ok := c.Data.([]interface{})
	if !ok {
		managerLogger.Log(logger.ERROR, "in AddAd data err")
		return
	}

	type AdAddFmt struct {
		Ad_id             string
		Order_id          string
		Adtype            string
		Duration          string
		Mimes             string
		Channel           string
		Cs                string
		Size              string
		Tmp_name          string
		Priority          string
		Thirdparty_url    string
		Landingpage_url   string
		Creative_url      string
		Html_snippet      string
		Exchange          []string // 在哪儿？
		Category          []string
		Category_product  []string
		Attribute         string
		Buyer_creative_id string
		Advertiser_name   string
		Active            string

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
			switch src.Adtype {
			case "banner":
				dst.Type = common.AdType(common.Banner)
			case "video":
				dst.Type = common.AdType(common.Video)
			default:
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
		common.GAdContainer.Add(&ad)
	}
}

func (c *Command) DelAd() {
	data, ok := c.Data.([]interface{})
	if !ok {
		managerLogger.Log(logger.ERROR, "in DelAd data err")
		return
	}
	for _, d := range data {
		if id, ok := d.(string); ok {
			common.GAdContainer.Del(id)
		} else {
			managerLogger.Log(logger.ERROR, "in DelAd, data array fmt err")
		}
	}
}

var managerLogger *logger.Log

func Init(path string) {
	managerLogger = logger.NewLog(path)
}
