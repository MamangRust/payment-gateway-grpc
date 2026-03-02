package topup_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type topupCachedResponseAll struct {
	Data  []*db.GetTopupsRow `json:"data"`
	Total *int               `json:"total_records"`
}

type topupCachedResponseByCard struct {
	Data  []*db.GetTopupsByCardNumberRow `json:"data"`
	Total *int                           `json:"total_records"`
}

type topupCachedResponseActive struct {
	Data  []*db.GetActiveTopupsRow `json:"data"`
	Total *int                     `json:"total_records"`
}

type topupCachedResponseTrashed struct {
	Data  []*db.GetTrashedTopupsRow `json:"data"`
	Total *int                      `json:"total_records"`
}

type topupQueryCache struct {
	store *cache.CacheStore
}

func NewTopupQueryCache(store *cache.CacheStore) TopupQueryCache {
	return &topupQueryCache{store: store}
}

func (c *topupQueryCache) GetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTopupsRow, *int, bool) {
	key := fmt.Sprintf(topupAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[topupCachedResponseAll](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupsCache(ctx context.Context, req *requests.FindAllTopups, data []*db.GetTopupsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTopupsRow{}
	}

	key := fmt.Sprintf(topupAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &topupCachedResponseAll{Data: data, Total: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*db.GetTopupsByCardNumberRow, *int, bool) {
	key := fmt.Sprintf(topupByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[topupCachedResponseByCard](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCacheTopupByCardCache(ctx context.Context, req *requests.FindAllTopupsByCardNumber, data []*db.GetTopupsByCardNumberRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTopupsByCardNumberRow{}
	}

	key := fmt.Sprintf(topupByCardCacheKey, req.CardNumber, req.Page, req.PageSize, req.Search)
	payload := &topupCachedResponseByCard{Data: data, Total: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetActiveTopupsRow, *int, bool) {
	key := fmt.Sprintf(topupActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[topupCachedResponseActive](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupActiveCache(ctx context.Context, req *requests.FindAllTopups, data []*db.GetActiveTopupsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetActiveTopupsRow{}
	}

	key := fmt.Sprintf(topupActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &topupCachedResponseActive{Data: data, Total: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTrashedTopupsRow, *int, bool) {
	key := fmt.Sprintf(topupTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[topupCachedResponseTrashed](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *topupQueryCache) SetCachedTopupTrashedCache(ctx context.Context, req *requests.FindAllTopups, data []*db.GetTrashedTopupsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTrashedTopupsRow{}
	}

	key := fmt.Sprintf(topupTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &topupCachedResponseTrashed{Data: data, Total: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *topupQueryCache) GetCachedTopupCache(ctx context.Context, id int) (*db.GetTopupByIDRow, bool) {
	key := fmt.Sprintf(topupByIdCacheKey, id)

	result, found := cache.GetFromCache[*db.GetTopupByIDRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupQueryCache) SetCachedTopupCache(ctx context.Context, data *db.GetTopupByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(topupByIdCacheKey, data.TopupID)
	cache.SetToCache(ctx, c.store, key, data, ttlDefault)
}
