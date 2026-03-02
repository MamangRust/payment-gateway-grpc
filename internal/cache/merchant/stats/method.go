package merchant_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type merchantStatsMethodCache struct {
	store *cache.CacheStore
}

func NewMerchantStatsMethodCache(store *cache.CacheStore) MerchantStatsMethodCache {
	return &merchantStatsMethodCache{store: store}
}

func (s *merchantStatsMethodCache) GetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsMerchantRow, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyPaymentMethodsMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsMethodCache) SetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int, data []*db.GetMonthlyPaymentMethodsMerchantRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodCacheKey, year)

	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *merchantStatsMethodCache) GetYearlyPaymentMethodMerchantCache(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodMerchantRow, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyPaymentMethodMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantStatsMethodCache) SetYearlyPaymentMethodMerchantCache(ctx context.Context, year int, data []*db.GetYearlyPaymentMethodMerchantRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodCacheKey, year)

	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
