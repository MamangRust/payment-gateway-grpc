package withdraw_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type withdrawStatsAmountCache struct {
	store *cache.CacheStore
}

func NewWithdrawStatsAmountCache(store *cache.CacheStore) WithdrawStatsAmountCache {
	return &withdrawStatsAmountCache{store: store}
}

func (w *withdrawStatsAmountCache) GetCachedMonthlyWithdraws(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawsRow, bool) {
	key := fmt.Sprintf(montWithdrawAmountKey, year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyWithdrawsRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsAmountCache) SetCachedMonthlyWithdraws(ctx context.Context, year int, data []*db.GetMonthlyWithdrawsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawAmountKey, year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsAmountCache) GetCachedYearlyWithdraws(ctx context.Context, year int) ([]*db.GetYearlyWithdrawsRow, bool) {
	key := fmt.Sprintf(yearWithdrawAmountKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyWithdrawsRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsAmountCache) SetCachedYearlyWithdraws(ctx context.Context, year int, data []*db.GetYearlyWithdrawsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawAmountKey, year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}
