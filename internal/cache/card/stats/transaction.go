package card_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsTransactionCache struct {
	store *cache.CacheStore
}

func NewCardStatsTransactionCache(store *cache.CacheStore) CardStatsTransactionCache {
	return &cardStatsTransactionCache{store: store}
}

func (c *cardStatsTransactionCache) GetMonthlyTransactionCache(ctx context.Context, year int) ([]*db.GetMonthlyTransactionAmountRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransactionAmount, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionAmountRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransactionCache) SetMonthlyTransactionCache(ctx context.Context, year int, data []*db.GetMonthlyTransactionAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransactionAmount, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTransactionCache) GetYearlyTransactionCache(ctx context.Context, year int) ([]*db.GetYearlyTransactionAmountRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransactionAmount, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionAmountRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransactionCache) SetYearlyTransactionCache(ctx context.Context, year int, data []*db.GetYearlyTransactionAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransactionAmount, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
