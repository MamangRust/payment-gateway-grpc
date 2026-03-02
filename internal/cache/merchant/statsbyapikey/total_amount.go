package merchant_stats_byapikey_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type merchantStatsTotalAmountByApiKeyCache struct {
	store *cache.CacheStore
}

func NewMerchantStatsTotalAmountByApiKeyCache(store *cache.CacheStore) MerchantStatsTotalAmountByApiKeyCache {
	return &merchantStatsTotalAmountByApiKeyCache{store: store}
}

func (m *merchantStatsTotalAmountByApiKeyCache) GetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetMonthlyTotalAmountByApikeyRow, bool) {
	key := fmt.Sprintf(merchantMonthlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalAmountByApikeyRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsTotalAmountByApiKeyCache) SetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*db.GetMonthlyTotalAmountByApikeyRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsTotalAmountByApiKeyCache) GetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetYearlyTotalAmountByApikeyRow, bool) {
	key := fmt.Sprintf(merchantYearlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTotalAmountByApikeyRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsTotalAmountByApiKeyCache) SetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*db.GetYearlyTotalAmountByApikeyRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyTotalAmountByApikeyCacheKey, req.Apikey, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
