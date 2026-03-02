package merchant_stats_byapikey_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type merchantStatsMethodByApiKeyCache struct {
	store *cache.CacheStore
}

func NewMerchantStatsMethodByApiKeyCache(store *cache.CacheStore) MerchantStatsMethodByApiKeyCache {
	return &merchantStatsMethodByApiKeyCache{store: store}
}

func (m *merchantStatsMethodByApiKeyCache) GetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetMonthlyPaymentMethodByApikeyRow, bool) {
	key := fmt.Sprintf(merchantMonthlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyPaymentMethodByApikeyRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsMethodByApiKeyCache) SetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*db.GetMonthlyPaymentMethodByApikeyRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsMethodByApiKeyCache) GetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetYearlyPaymentMethodByApikeyRow, bool) {
	key := fmt.Sprintf(merchantYearlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyPaymentMethodByApikeyRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsMethodByApiKeyCache) SetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*db.GetYearlyPaymentMethodByApikeyRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyPaymentMethodByApikeyCacheKey, req.Apikey, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
