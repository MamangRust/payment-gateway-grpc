package merchant_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type merchantStatsTotalAmountCache struct {
	store *cache.CacheStore
}

func NewMerchantStatsTotalAmountCache(store *cache.CacheStore) MerchantStatsTotalAmountCache {
	return &merchantStatsTotalAmountCache{store: store}
}

func (s *merchantStatsTotalAmountCache) GetMonthlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*db.GetMonthlyTotalAmountMerchantRow, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalAmountMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsTotalAmountCache) SetMonthlyTotalAmountMerchantCache(ctx context.Context, year int, data []*db.GetMonthlyTotalAmountMerchantRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyTotalAmountCacheKey, year)

	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *merchantStatsTotalAmountCache) GetYearlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*db.GetYearlyTotalAmountMerchantRow, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTotalAmountMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsTotalAmountCache) SetYearlyTotalAmountMerchantCache(ctx context.Context, year int, data []*db.GetYearlyTotalAmountMerchantRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyTotalAmountCacheKey, year)

	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
