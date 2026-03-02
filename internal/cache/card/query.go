package card_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardCachedResponse struct {
	Data         []*db.GetCardsRow `json:"data"`
	TotalRecords *int              `json:"total_records"`
}

type cardCachedResponseActive struct {
	Data         []*db.GetActiveCardsWithCountRow `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type cardCachedResponseTrashed struct {
	Data         []*db.GetTrashedCardsWithCountRow `json:"data"`
	TotalRecords *int                              `json:"total_records"`
}

type cardQueryCache struct {
	store *cache.CacheStore
}

func NewCardQueryCache(store *cache.CacheStore) CardQueryCache {
	return &cardQueryCache{store: store}
}

func (c *cardQueryCache) GetByIdCache(ctx context.Context, cardID int) (*db.GetCardByIDRow, bool) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)

	result, found := cache.GetFromCache[*db.GetCardByIDRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetByUserIDCache(ctx context.Context, userID int) (*db.GetCardByUserIDRow, bool) {
	key := fmt.Sprintf(cardByUserIdCacheKey, userID)

	result, found := cache.GetFromCache[*db.GetCardByUserIDRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetByCardNumberCache(ctx context.Context, cardNumber string) (*db.GetCardByCardNumberRow, bool) {
	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)

	result, found := cache.GetFromCache[*db.GetCardByCardNumberRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardQueryCache) GetFindAllCache(ctx context.Context, req *requests.FindAllCards) ([]*db.GetCardsRow, *int, bool) {
	key := fmt.Sprintf(cardAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[cardCachedResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) GetByActiveCache(ctx context.Context, req *requests.FindAllCards) ([]*db.GetActiveCardsWithCountRow, *int, bool) {
	key := fmt.Sprintf(cardActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[cardCachedResponseActive](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) GetByTrashedCache(ctx context.Context, req *requests.FindAllCards) ([]*db.GetTrashedCardsWithCountRow, *int, bool) {
	key := fmt.Sprintf(cardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[cardCachedResponseTrashed](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (c *cardQueryCache) SetByIdCache(ctx context.Context, cardID int, data *db.GetCardByIDRow) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)
	cache.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetByUserIDCache(ctx context.Context, userID int, data *db.GetCardByUserIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardByUserIdCacheKey, userID)
	cache.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetByCardNumberCache(ctx context.Context, cardNumber string, data *db.GetCardByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	cache.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *cardQueryCache) SetFindAllCache(ctx context.Context, req *requests.FindAllCards, data []*db.GetCardsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetCardsRow{}
	}

	payload := &cardCachedResponse{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardAllCacheKey, req.Page, req.PageSize, req.Search)
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) SetByActiveCache(ctx context.Context, req *requests.FindAllCards, data []*db.GetActiveCardsWithCountRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetActiveCardsWithCountRow{}
	}

	payload := &cardCachedResponseActive{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardActiveCacheKey, req.Page, req.PageSize, req.Search)
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) SetByTrashedCache(ctx context.Context, req *requests.FindAllCards, data []*db.GetTrashedCardsWithCountRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetTrashedCardsWithCountRow{}
	}

	payload := &cardCachedResponseTrashed{Data: data, TotalRecords: total}

	key := fmt.Sprintf(cardTrashedCacheKey, req.Page, req.PageSize, req.Search)
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *cardQueryCache) DeleteByIdCache(ctx context.Context, cardID int) {
	key := fmt.Sprintf(cardByIdCacheKey, cardID)
	cache.DeleteFromCache(ctx, c.store, key)
}

func (c *cardQueryCache) DeleteByUserIDCache(ctx context.Context, userID int) {
	key := fmt.Sprintf(cardByUserIdCacheKey, userID)
	cache.DeleteFromCache(ctx, c.store, key)
}

func (c *cardQueryCache) DeleteByCardNumberCache(ctx context.Context, cardNumber string) {
	key := fmt.Sprintf(cardByCardNumCacheKey, cardNumber)
	cache.DeleteFromCache(ctx, c.store, key)
}
