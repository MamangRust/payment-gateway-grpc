package card_dashboard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
)

type cardDashboardCache struct {
	store *cache.CacheStore
}

func NewCardDashboardCache(store *cache.CacheStore) CardDashboardTotalCache {
	return &cardDashboardCache{store: store}
}

func (c *cardDashboardCache) GetDashboardCardCache(ctx context.Context) (*response.DashboardCard, bool) {
	result, found := cache.GetFromCache[*response.DashboardCard](ctx, c.store, cacheKeyDashboardDefault)

	if !found || result == nil {
		return nil, false
	}

	return *result, true

}

func (c *cardDashboardCache) SetDashboardCardCache(ctx context.Context, data *response.DashboardCard) {
	if data == nil {
		return
	}

	cache.SetToCache(ctx, c.store, cacheKeyDashboardDefault, data, ttlDashboardDefault)
}

func (c *cardDashboardCache) DeleteDashboardCardCache(ctx context.Context) {
	cache.DeleteFromCache(ctx, c.store, cacheKeyDashboardDefault)
}
