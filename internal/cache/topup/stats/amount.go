package topup_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type topupStatsAmountCache struct {
	store *cache.CacheStore
}

func NewTopupStatsAmountCache(store *cache.CacheStore) TopupStatsAmountCache {
	return &topupStatsAmountCache{store: store}
}

func (c *topupStatsAmountCache) GetMonthlyTopupAmountsCache(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountsRow, bool) {
	key := fmt.Sprintf(monthTopupAmountCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTopupAmountsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatsAmountCache) SetMonthlyTopupAmountsCache(ctx context.Context, year int, data []*db.GetMonthlyTopupAmountsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *topupStatsAmountCache) GetYearlyTopupAmountsCache(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountsRow, bool) {
	key := fmt.Sprintf(yearTopupAmountCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupAmountsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatsAmountCache) SetYearlyTopupAmountsCache(ctx context.Context, year int, data []*db.GetYearlyTopupAmountsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
