package merchant_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type merchantStatsAmountCache struct {
	store *cache.CacheStore
}

func NewMerchantStatsAmountCache(store *cache.CacheStore) MerchantStatsAmountCache {
	return &merchantStatsAmountCache{store: store}
}

func (s *merchantStatsAmountCache) GetMonthlyAmountMerchantCache(ctx context.Context, year int) ([]*db.GetMonthlyAmountMerchantRow, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsAmountCache) SetMonthlyAmountMerchantCache(ctx context.Context, year int, data []*db.GetMonthlyAmountMerchantRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountCacheKey, year)

	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *merchantStatsAmountCache) GetYearlyAmountMerchantCache(ctx context.Context, year int) ([]*db.GetYearlyAmountMerchantRow, bool) {
	key := fmt.Sprintf(MerchantYearlyAmountCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsAmountCache) SetYearlyAmountMerchantCache(ctx context.Context, year int, data []*db.GetYearlyAmountMerchantRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(MerchantYearlyAmountCacheKey, year)

	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
