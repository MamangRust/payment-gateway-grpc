package card_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
	"fmt"
)

type cardStatsWithdrawByCardCache struct {
	store *cache.CacheStore
}

func NewCardStatsWithdrawByCardCache(store *cache.CacheStore) CardStatsWithdrawByCardCache {
	return &cardStatsWithdrawByCardCache{store: store}
}

func (c *cardStatsWithdrawByCardCache) GetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) (*response.ApiResponseMonthlyAmount, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawByCard, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseMonthlyAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (c *cardStatsWithdrawByCardCache) SetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data *response.ApiResponseMonthlyAmount) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(cacheKeyMonthlyWithdrawByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, data, expirationCardStatistic)
}

func (c *cardStatsWithdrawByCardCache) GetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) (*response.ApiResponseYearlyAmount, bool) {
	key := fmt.Sprintf(cacheKeyYearlyWithdrawByCard, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseYearlyAmount](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (c *cardStatsWithdrawByCardCache) SetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data *response.ApiResponseYearlyAmount) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(cacheKeyYearlyWithdrawByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, data, expirationCardStatistic)
}
