package card_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsBalanceByCardCache struct {
	store *cache.CacheStore
}

func NewCardStatsBalanceByCardCache(store *cache.CacheStore) CardStatsBalanceByCardCache {
	return &cardStatsBalanceByCardCache{store: store}
}

func (c *cardStatsBalanceByCardCache) GetMonthlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyBalancesByCardNumberRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyBalanceByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyBalancesByCardNumberRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsBalanceByCardCache) SetMonthlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetMonthlyBalancesByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyBalanceByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsBalanceByCardCache) GetYearlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyBalancesByCardNumberRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyBalanceByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyBalancesByCardNumberRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsBalanceByCardCache) SetYearlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetYearlyBalancesByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyBalanceByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}
