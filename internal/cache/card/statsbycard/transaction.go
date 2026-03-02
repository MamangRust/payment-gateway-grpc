package card_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsTransactionByCardCache struct {
	store *cache.CacheStore
}

func NewCardStatsTransactionByCardCache(store *cache.CacheStore) CardStatsTransactionByCardCache {
	return &cardStatsTransactionByCardCache{store: store}
}

func (c *cardStatsTransactionByCardCache) GetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransactionAmountByCardNumberRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTxnByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionAmountByCardNumberRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransactionByCardCache) SetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetMonthlyTransactionAmountByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTxnByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTransactionByCardCache) GetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransactionAmountByCardNumberRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTxnByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionAmountByCardNumberRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransactionByCardCache) SetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetYearlyTransactionAmountByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTxnByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}
