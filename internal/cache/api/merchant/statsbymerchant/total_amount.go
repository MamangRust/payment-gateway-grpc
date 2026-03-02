package merchant_stats_bymerchant_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
	"fmt"
)

type merchantStatsTotalAmountByMerchant struct {
	store *cache.CacheStore
}

func NewMerchantStatsTotalAmountByMerchantCache(store *cache.CacheStore) MerchantStatsTotalAmountByMerchantCache {
	return &merchantStatsTotalAmountByMerchant{store: store}
}

func (m *merchantStatsTotalAmountByMerchant) GetMonthlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) (*response.ApiResponseMerchantMonthlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseMerchantMonthlyTotalAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (m *merchantStatsTotalAmountByMerchant) SetMonthlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant, data *response.ApiResponseMerchantMonthlyTotalAmount) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantMonthlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantStatsTotalAmountByMerchant) GetYearlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) (*response.ApiResponseMerchantYearlyTotalAmount, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseMerchantYearlyTotalAmount](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (m *merchantStatsTotalAmountByMerchant) SetYearlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant, data *response.ApiResponseMerchantYearlyTotalAmount) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantYearlyTotalAmountByMerchantCacheKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
