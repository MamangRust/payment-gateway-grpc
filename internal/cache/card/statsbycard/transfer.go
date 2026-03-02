package card_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsTransferByCardCache struct {
	store *cache.CacheStore
}

func NewCardStatsTransferByCardCache(store *cache.CacheStore) CardStatsTransferByCardCache {
	return &cardStatsTransferByCardCache{store: store}
}

func (c *cardStatsTransferByCardCache) GetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountBySenderRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlySenderByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransferAmountBySenderRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferByCardCache) SetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetMonthlyTransferAmountBySenderRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlySenderByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTransferByCardCache) GetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountBySenderRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlySenderByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransferAmountBySenderRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferByCardCache) SetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetYearlyTransferAmountBySenderRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlySenderByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTransferByCardCache) GetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountByReceiverRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyReceiverByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransferAmountByReceiverRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferByCardCache) SetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetMonthlyTransferAmountByReceiverRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyReceiverByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}

func (c *cardStatsTransferByCardCache) GetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountByReceiverRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyReceiverByCard, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransferAmountByReceiverRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferByCardCache) SetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*db.GetYearlyTransferAmountByReceiverRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyReceiverByCard, req.CardNumber, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, expirationCardStatistic)
}
