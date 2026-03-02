package topup_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
	"fmt"
)

type topupStatsMethodByCardCache struct {
	store *cache.CacheStore
}

func NewTopupStatsMethodByCardCache(store *cache.CacheStore) TopupStatsMethodByCardCache {
	return &topupStatsMethodByCardCache{store: store}
}

func (s *topupStatsMethodByCardCache) GetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) (*response.ApiResponseTopupMonthMethod, bool) {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseTopupMonthMethod](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *topupStatsMethodByCardCache) SetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data *response.ApiResponseTopupMonthMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *topupStatsMethodByCardCache) GetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) (*response.ApiResponseTopupYearMethod, bool) {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseTopupYearMethod](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *topupStatsMethodByCardCache) SetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data *response.ApiResponseTopupYearMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
