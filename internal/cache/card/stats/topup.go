package card_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsTopupCache struct {
	store *cache.CacheStore
}

func NewCardStatsTopupCache(store *cache.CacheStore) CardStatsTopupCache {
	return &cardStatsTopupCache{store: store}
}

func (c *cardStatsTopupCache) GetMonthlyTopupCache(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTopupAmount, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTopupAmountRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTopupCache) SetMonthlyTopupCache(ctx context.Context, year int, data []*db.GetMonthlyTopupAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTopupAmount, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTopupCache) GetYearlyTopupCache(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTopupAmount, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupAmountRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTopupCache) SetYearlyTopupCache(ctx context.Context, year int, data []*db.GetYearlyTopupAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTopupAmount, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
