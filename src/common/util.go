package common

import (
	"errors"
	"sync"
)

var WinUrl string
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

func (c *AdContainer) Add(ad *Ad) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if _, ok := c.Idmap[ad.Id]; ok {
		return
	}
	c.Idmap[ad.Id] = len(c.Ads) + 1
	c.Ads = append(c.Ads, *ad)
}

func (c *AdContainer) Del(id string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	idx, ok := c.Idmap[id]
	if !ok {
		return
	}
	defer delete(c.Idmap, id)
	if idx >= len(c.Ads) {
		return
	}
	idx2 := len(c.Ads) - 1
	defer func() {
		copy(c.Ads[:], c.Ads[0:idx2])
	}()
	if idx != idx2 {
		c.Ads[idx] = c.Ads[idx2]
		c.Idmap[c.Ads[idx2].Id] = idx
	}
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

func (c *OrderContainer) Add(order *Order) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if _, ok := c.Idmap[order.Id]; ok {
		return
	}
	c.Idmap[order.Id] = len(c.Orders) + 1
	c.Orders = append(c.Orders, *order)
}

func (c *OrderContainer) Del(id string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	idx, ok := c.Idmap[id]
	if !ok {
		return
	}
	defer delete(c.Idmap, id)
	if idx >= len(c.Orders) {
		return
	}
	idx2 := len(c.Orders) - 1
	defer func() {
		copy(c.Orders[:], c.Orders[0:idx2])
	}()
	if idx != idx2 {
		c.Orders[idx] = c.Orders[idx2]
		c.Idmap[c.Orders[idx2].Id] = idx
	}
}

func (c *OrderContainer) find(id string) int {
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
