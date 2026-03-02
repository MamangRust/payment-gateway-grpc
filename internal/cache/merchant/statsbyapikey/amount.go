package merchant_stats_byapikey_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type merchantStatsAmountByApiKeyCache struct {
	store *cache.CacheStore
}

func NewMerchantStatsAmountByApiKeyCache(store *cache.CacheStore) MerchantStatsAmountByApiKeyCache {
	return &merchantStatsAmountByApiKeyCache{store: store}
}

func (m *merchantStatsAmountByApiKeyCache) GetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetMonthlyAmountByApikeyRow, bool) {
	key := fmt.Sprintf(merchantMonthlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountByApikeyRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsAmountByApiKeyCache) SetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*db.GetMonthlyAmountByApikeyRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantMonthlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *merchantStatsAmountByApiKeyCache) GetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetYearlyAmountByApikeyRow, bool) {
	key := fmt.Sprintf(merchantYearlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountByApikeyRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantStatsAmountByApiKeyCache) SetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*db.GetYearlyAmountByApikeyRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantYearlyAmountByApikeyCacheKey, req.Apikey, req.Year)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
