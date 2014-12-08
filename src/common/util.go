package common

import (
	"errors"
	"fmt"
	"sync"
)

//var WinUrl string (这个应该是各个handler自己填)
var GAdContainer *AdContainer = NewAdContainer()
var GOrderContainer *OrderContainer = NewOrderContainer()

type AdContainer struct {
	Idmap map[string]int
	Ads   []Ad
	Lock  sync.RWMutex
}

func NewAdContainer() *AdContainer {
	return &AdContainer{
		Idmap: make(map[string]int),
		Ads:   make([]Ad, 0, 1024),
	}
}

func (c *AdContainer) Add(ad *Ad) int {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if _, ok := c.Idmap[ad.Id]; ok {
		return len(c.Ads)
	}
	c.Idmap[ad.Id] = len(c.Ads)
	c.Ads = append(c.Ads, *ad)
	return len(c.Ads)
}

func (c *AdContainer) Del(id string) int {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	idx, ok := c.Idmap[id]
	if !ok {
		return len(c.Ads)
	}
	defer delete(c.Idmap, id)
	if idx >= len(c.Ads) {
		return len(c.Ads)
	}
	idx2 := len(c.Ads) - 1
	defer func() {
		copy(c.Ads[:], c.Ads[0:idx2])
	}()
	if idx != idx2 {
		c.Ads[idx] = c.Ads[idx2]
		c.Idmap[c.Ads[idx2].Id] = idx
	}
	return len(c.Ads)
}

func (c *AdContainer) Find(id string) (ad Ad, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	if v, ok := c.Idmap[id]; !ok {
		err = errors.New("No Ad Found: " + id)
		return
	} else {
		ad = c.Ads[v]
		return
	}
}

type OrderContainer struct {
	Idmap  map[string]int
	Orders []Order
	Lock   sync.RWMutex
}

func NewOrderContainer() *OrderContainer {
	return &OrderContainer{
		Idmap:  make(map[string]int),
		Orders: make([]Order, 0, 1024),
	}
}

func (c *OrderContainer) SetCost(orderId string, cost int) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	for i := 0; i < len(c.Orders); i++ {
		if c.Orders[i].Id == orderId {
			c.Orders[i].CountCost = cost
			return
		}
	}
	fmt.Println("SetCost order: ", orderId, "not exist")
}

func (c *OrderContainer) Add(order *Order) int {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if _, ok := c.Idmap[order.Id]; ok {
		return len(c.Orders)
	}
	c.Idmap[order.Id] = len(c.Orders)
	c.Orders = append(c.Orders, *order)
	return len(c.Orders)
}

func (c *OrderContainer) Del(id string) int {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	idx, ok := c.Idmap[id]
	if !ok {
		return len(c.Orders)
	}
	defer delete(c.Idmap, id)
	if idx >= len(c.Orders) {
		return len(c.Orders)
	}
	idx2 := len(c.Orders) - 1
	defer func() {
		copy(c.Orders[:], c.Orders[0:idx2])
	}()
	if idx != idx2 {
		c.Orders[idx] = c.Orders[idx2]
		c.Idmap[c.Orders[idx2].Id] = idx
	}
	return len(c.Orders)
}

func (c *OrderContainer) find(id string) int {
	//fmt.Printf("Idmap = %#v\n", c.Idmap)
	if idx, ok := c.Idmap[id]; ok {
		return idx
	}
	return -1
}

func (c *OrderContainer) Find(id string) (order Order, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	if idx := c.find(id); idx != -1 {
		order = c.Orders[idx]
	} else {
		err = errors.New("No order Found: " + id)
	}
	return
}

func (c *OrderContainer) FindPrice(id string, adx Adx) (int, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	if idx := c.find(id); idx == -1 {
		return 0, errors.New("No order found: " + id)
	} else if len(c.Orders[idx].MaxPrice) <= idx {
		return 0, errors.New("order " + id + " price slice error")
	} else {
		return c.Orders[idx].MaxPrice[int(adx)], nil
	}
}

func (c *OrderContainer) FeeEnough(id string) bool {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	if idx := c.find(id); idx == -1 {
		fmt.Println("cannot found order id: ", id)
		return false
	} else if c.Orders[idx].CountCost >= 0 {
		return true
	} else {
		fmt.Println("order ", id, "CountCost not enough: ", c.Orders[idx].CountCost)
		return false
	}
}

func (c *OrderContainer) FeeDecr(id string, cost int) {
	if cost <= 0 {
		fmt.Println("FeeDecr: cost shoud be a post number")
		return
	}
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	if idx := c.find(id); idx == -1 {
		return
	} else {
		c.Orders[idx].CountCost -= cost
	}
}
