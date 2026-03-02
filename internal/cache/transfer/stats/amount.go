package transfer_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transferStatsAmountCache struct {
	store *cache.CacheStore
}

func NewTransferStatsAmountCache(store *cache.CacheStore) TransferStatsAmountCache {
	return &transferStatsAmountCache{store: store}
}

func (t *transferStatsAmountCache) GetCachedMonthTransferAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountsRow, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountKey, year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTransferAmountsRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsAmountCache) SetCachedMonthTransferAmounts(ctx context.Context, year int, data []*db.GetMonthlyTransferAmountsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsAmountCache) GetCachedYearlyTransferAmounts(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountsRow, bool) {
	key := fmt.Sprintf(transferYearTransferAmountKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransferAmountsRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsAmountCache) SetCachedYearlyTransferAmounts(ctx context.Context, year int, data []*db.GetYearlyTransferAmountsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
