package transfer_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transferCacheResponseAll struct {
	Data         []*db.GetTransfersRow `json:"data"`
	TotalRecords *int                  `json:"total_records"`
}

type transferCacheResponseActive struct {
	Data         []*db.GetActiveTransfersRow `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type transferCacheResponseTrashed struct {
	Data         []*db.GetTrashedTransfersRow `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type transferQueryCache struct {
	store *cache.CacheStore
}

func NewTransferQueryCache(store *cache.CacheStore) TransferQueryCache {
	return &transferQueryCache{store: store}
}

func (c *transferQueryCache) GetCachedTransfersCache(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTransfersRow, *int, bool) {
	key := fmt.Sprintf(transferAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[transferCacheResponseAll](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (c *transferQueryCache) SetCachedTransfersCache(ctx context.Context, req *requests.FindAllTranfers, data []*db.GetTransfersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTransfersRow{}
	}

	key := fmt.Sprintf(transferAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transferCacheResponseAll{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetActiveTransfersRow, *int, bool) {
	key := fmt.Sprintf(transferActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[transferCacheResponseActive](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *transferQueryCache) SetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTranfers, data []*db.GetActiveTransfersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetActiveTransfersRow{}
	}

	key := fmt.Sprintf(transferActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transferCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTrashedTransfersRow, *int, bool) {
	key := fmt.Sprintf(transferTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[transferCacheResponseTrashed](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *transferQueryCache) SetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTranfers, data []*db.GetTrashedTransfersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTrashedTransfersRow{}
	}

	key := fmt.Sprintf(transferTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transferCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferCache(ctx context.Context, id int) (*db.GetTransferByIDRow, bool) {
	key := fmt.Sprintf(transferByIdCacheKey, id)
	result, found := cache.GetFromCache[*db.GetTransferByIDRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *transferQueryCache) SetCachedTransferCache(ctx context.Context, data *db.GetTransferByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByIdCacheKey, data.TransferID)
	cache.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferByFrom(ctx context.Context, from string) ([]*db.GetTransfersBySourceCardRow, bool) {
	key := fmt.Sprintf(transferByFromCacheKey, from)
	result, found := cache.GetFromCache[[]*db.GetTransfersBySourceCardRow](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *transferQueryCache) SetCachedTransferByFrom(ctx context.Context, from string, data []*db.GetTransfersBySourceCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByFromCacheKey, from)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *transferQueryCache) GetCachedTransferByTo(ctx context.Context, to string) ([]*db.GetTransfersByDestinationCardRow, bool) {
	key := fmt.Sprintf(transferByToCacheKey, to)

	result, found := cache.GetFromCache[[]*db.GetTransfersByDestinationCardRow](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *transferQueryCache) SetCachedTransferByTo(ctx context.Context, to string, data []*db.GetTransfersByDestinationCardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferByToCacheKey, to)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
