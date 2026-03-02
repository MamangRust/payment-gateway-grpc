package transaction_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transactionStatsAmountCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsAmountCache(store *cache.CacheStore) TransactionStatsAmountCache {
	return &transactionStatsAmountCache{store: store}
}

func (t *transactionStatsAmountCache) GetMonthlyAmountsCache(ctx context.Context, year int) ([]*db.GetMonthlyAmountsRow, bool) {
	key := fmt.Sprintf(monthTransactionAmountCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountsRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsAmountCache) SetMonthlyAmountsCache(ctx context.Context, year int, data []*db.GetMonthlyAmountsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionAmountCacheKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsAmountCache) GetYearlyAmountsCache(ctx context.Context, year int) ([]*db.GetYearlyAmountsRow, bool) {
	key := fmt.Sprintf(yearTransactionAmountCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountsRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsAmountCache) SetYearlyAmountsCache(ctx context.Context, year int, data []*db.GetYearlyAmountsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionAmountCacheKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
