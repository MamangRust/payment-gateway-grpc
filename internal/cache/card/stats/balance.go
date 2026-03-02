package card_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsBalanceCache struct {
	store *cache.CacheStore
}

func NewCardStatsBalanceCache(store *cache.CacheStore) CardStatsBalanceCache {
	return &cardStatsBalanceCache{store: store}
}

func (c *cardStatsBalanceCache) GetMonthlyBalanceCache(ctx context.Context, year int) ([]*db.GetMonthlyBalancesRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyBalance, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyBalancesRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsBalanceCache) SetMonthlyBalanceCache(ctx context.Context, year int, data []*db.GetMonthlyBalancesRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyBalance, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsBalanceCache) GetYearlyBalanceCache(ctx context.Context, year int) ([]*db.GetYearlyBalancesRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyBalance, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyBalancesRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsBalanceCache) SetYearlyBalanceCache(ctx context.Context, year int, data []*db.GetYearlyBalancesRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyBalance, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
