package cachestore

import (
	"errors"

	"github.com/kamva/hexa"
	hcache "github.com/kamva/hexa-cache"
	"github.com/kamva/hexa/hlog"
	"space.org/space/internal/model"
	"space.org/space/internal/registry/services"
)

type CacheStore struct {
	model.Store
	stores *storesList
}

type storesList struct {
	//user       model.UserStore

	// Place other stores here
}

//func (c *CacheStore) User() model.UserStore {
//	return c.stores.user
//}

// Add other methods to return cache store here.

func New(r hexa.ServiceRegistry, next model.Store) *CacheStore {
	cs := &CacheStore{Store: next}
	cp := services.New(r).CacheProvider()
	_ = cp

	cs.stores = &storesList{
		//user:       &userCacheStore{UserStore: next.User(), rootCache: cs, cache: cp.Cache("user")},

		// Add other cache implementations here.
	}

	return cs
}

func handleCacheErr(err error) error {
	if err == nil {
		return nil
	}

	if !errors.Is(err, hcache.ErrKeyNotFound) {
		hlog.Error("can not fetch data from cache", hlog.Err(err))
		return nil
	}

	return err
}

func logErr(err error) {
	if err != nil {
		hlog.Error("can not set/remove data on the cache server", hlog.Err(err))
	}
}

var _ model.Store = &CacheStore{}
