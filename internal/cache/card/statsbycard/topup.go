package card_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsTopupByCardCache struct {
	store *cache.CacheStore
}

func NewCardStatsTopupByCardCache(store *cache.CacheStore) CardStatsTopupByCardCache {
	return &cardStatsTopupByCardCache{store: store}
}

func (c *cardStatsTopupByCardCache) GetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTopupAmountByCardNumberRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTopupByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTopupAmountByCardNumberRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTopupByCardCache) SetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetMonthlyTopupAmountByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTopupByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTopupByCardCache) GetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTopupAmountByCardNumberRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTopupByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupAmountByCardNumberRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTopupByCardCache) SetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetYearlyTopupAmountByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTopupByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}
