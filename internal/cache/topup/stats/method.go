package topup_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type topupStatsMethodCache struct {
	store *cache.CacheStore
}

func NewTopupStatsMethodCache(store *cache.CacheStore) TopupStatsMethodCache {
	return &topupStatsMethodCache{store: store}
}

func (c *topupStatsMethodCache) GetMonthlyTopupMethodsCache(ctx context.Context, year int) ([]*db.GetMonthlyTopupMethodsRow, bool) {
	key := fmt.Sprintf(monthTopupMethodCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTopupMethodsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *topupStatsMethodCache) SetMonthlyTopupMethodsCache(ctx context.Context, year int, data []*db.GetMonthlyTopupMethodsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *topupStatsMethodCache) GetYearlyTopupMethodsCache(ctx context.Context, year int) ([]*db.GetYearlyTopupMethodsRow, bool) {
	key := fmt.Sprintf(yearTopupMethodCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupMethodsRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *topupStatsMethodCache) SetYearlyTopupMethodsCache(ctx context.Context, year int, data []*db.GetYearlyTopupMethodsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
