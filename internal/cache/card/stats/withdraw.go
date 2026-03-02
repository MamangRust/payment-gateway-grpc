package card_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsWithdrawCache struct {
	store *cache.CacheStore
}

func NewCardStatsWithdrawCache(store *cache.CacheStore) CardStatsWithdrawCache {
	return &cardStatsWithdrawCache{store: store}
}

func (c *cardStatsWithdrawCache) GetMonthlyWithdrawCache(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawAmountRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawAmount, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyWithdrawAmountRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsWithdrawCache) SetMonthlyWithdrawCache(ctx context.Context, year int, data []*db.GetMonthlyWithdrawAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyWithdrawAmount, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsWithdrawCache) GetYearlyWithdrawCache(ctx context.Context, year int) ([]*db.GetYearlyWithdrawAmountRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyWithdrawAmount, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyWithdrawAmountRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsWithdrawCache) SetYearlyWithdrawCache(ctx context.Context, year int, data []*db.GetYearlyWithdrawAmountRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyWithdrawAmount, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
