package card_dashboard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
	"fmt"
)

type cardDashboardByCardNumberCache struct {
	store *cache.CacheStore
}

func NewCardDashboardByCardNumberCache(store *cache.CacheStore) CardDashboardByCardNumberCache {
	return &cardDashboardByCardNumberCache{store: store}
}

func (c *cardDashboardByCardNumberCache) GetDashboardCardCardNumberCache(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, bool) {
	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	result, found := cache.GetFromCache[*response.DashboardCardCardNumber](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardDashboardByCardNumberCache) SetDashboardCardCardNumberCache(ctx context.Context, cardNumber string, data *response.DashboardCardCardNumber) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	cache.SetToCache(ctx, c.store, key, data, ttlDashboardDefault)
}

func (c *cardDashboardByCardNumberCache) DeleteDashboardCardCardNumberCache(ctx context.Context, cardNumber string) {
	key := fmt.Sprintf(cacheKeyDashboardCardNumber, cardNumber)
	cache.DeleteFromCache(ctx, c.store, key)
}
