package hotels

import (
	"context"
	"fmt"
	"github.com/karlseguin/ccache"
	hotelsDAO "hotels-api/dao/hotels"
)

type CacheConfig struct {
	MaxSize      int64
	ItemsToPrune uint32
}

type Cache struct {
	client *ccache.Cache
}

const (
	keyFormat = "Hotel:%d"
)

func NewCache(config CacheConfig) Cache {
	client := ccache.New(ccache.Configure[int]().MaxSize(config.MaxSize).ItemsToPrune(config.ItemsToPrune))
	return Cache{
		client: client,
	}
}

func (repository Cache) GetHotelByID(ctx context.Context, id int64) (hotelsDAO.Hotel, error) {
	key := fmt.Sprintf(keyFormat, id)
	item := repository.client.Get(key)
	if item.Expired() {
		return hotelsDAO.Hotel{}, fmt.Errorf("item %s is expired", key)
	}
	hotel, ok := item.Value().(hotelsDAO.Hotel)
	if !ok {
		return hotelsDAO.Hotel{}, fmt.Errorf("error converting item %s to DAO", key)
	}
	return hotel, nil
}
