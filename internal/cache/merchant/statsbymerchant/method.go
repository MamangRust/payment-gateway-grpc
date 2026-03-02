package merchant_stats_bymerchant_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type merchantStatsMethodByMerchant struct {
	store *cache.CacheStore
}

func NewMerchantStatsMethodByMerchantCache(store *cache.CacheStore) MerchantStatsMethodByMerchantCache {
	return &merchantStatsMethodByMerchant{store: store}
}

func (m *merchantStatsMethodByMerchant) GetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetMonthlyPaymentMethodByMerchantsRow, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyPaymentMethodByMerchantsRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsMethodByMerchant) SetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*db.GetMonthlyPaymentMethodByMerchantsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsMethodByMerchant) GetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetYearlyPaymentMethodByMerchantsRow, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyPaymentMethodByMerchantsRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsMethodByMerchant) SetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*db.GetYearlyPaymentMethodByMerchantsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodByMerchantCacheKey, req.MerchantID, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
