package cache

import (
	"encoding/json"
	"errors"
	"github.com/patrickmn/go-cache"
	"test0/internal/db"
	"test0/internal/models"
	"time"
)

type CacheMemory struct {
	Cache *cache.Cache
}

func NewCache(timeExpiration, timeCleanUp time.Duration) *CacheMemory {
	var c CacheMemory
	c.Cache = cache.New(timeExpiration, timeCleanUp)
	return &c
}

func (c *CacheMemory) Put(id string, order models.Order) {
	c.Cache.Set(id, order, cache.DefaultExpiration)
}

func (c *CacheMemory) GetOrder(id string) (json.RawMessage, error) {
	item, found := c.Cache.Get(id)
	if found {
		order, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		return order, nil
	}
	return nil, errors.New("Not found order :(")
}

func (c *CacheMemory) Restore() error {
	orders, err := db.GetAllOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		c.Put(order.OrderUID, order)
	}

	return nil
}
