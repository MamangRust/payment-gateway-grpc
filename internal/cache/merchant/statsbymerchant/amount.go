package merchant_stats_bymerchant_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type merchantStatsAmountByMerchant struct {
	store *cache.CacheStore
}

func NewMerchantStatsAmountByMerchantCache(store *cache.CacheStore) MerchantStatsAmountByMerchantCache {
	return &merchantStatsAmountByMerchant{store: store}
}

func (m *merchantStatsAmountByMerchant) GetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetMonthlyAmountByMerchantsRow, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountByMerchantsRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsAmountByMerchant) SetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*db.GetMonthlyAmountByMerchantsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsAmountByMerchant) GetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetYearlyAmountByMerchantsRow, bool) {
	key := fmt.Sprintf(merchantYearlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountByMerchantsRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsAmountByMerchant) SetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*db.GetYearlyAmountByMerchantsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyAmountByMerchantCacheKey, req.MerchantID, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
