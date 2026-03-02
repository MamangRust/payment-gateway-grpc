package card_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type cardStatsTransferCache struct {
	store *cache.CacheStore
}

func NewCardStatsTransferCache(store *cache.CacheStore) CardStatsTransferCache {
	return &cardStatsTransferCache{store: store}
}

func (c *cardStatsTransferCache) GetMonthlyTransferSenderCache(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountSenderRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransferSender, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransferAmountSenderRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferCache) SetMonthlyTransferSenderCache(ctx context.Context, year int, data []*db.GetMonthlyTransferAmountSenderRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransferSender, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTransferCache) GetYearlyTransferSenderCache(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountSenderRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransferSender, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransferAmountSenderRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferCache) SetYearlyTransferSenderCache(ctx context.Context, year int, data []*db.GetYearlyTransferAmountSenderRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransferSender, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTransferCache) GetMonthlyTransferReceiverCache(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountReceiverRow, bool) {
	key := fmt.Sprintf(cacheKeyMonthlyTransferReceiver, year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTransferAmountReceiverRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferCache) SetMonthlyTransferReceiverCache(ctx context.Context, year int, data []*db.GetMonthlyTransferAmountReceiverRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyMonthlyTransferReceiver, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}

func (c *cardStatsTransferCache) GetYearlyTransferReceiverCache(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountReceiverRow, bool) {
	key := fmt.Sprintf(cacheKeyYearlyTransferReceiver, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransferAmountReceiverRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *cardStatsTransferCache) SetYearlyTransferReceiverCache(ctx context.Context, year int, data []*db.GetYearlyTransferAmountReceiverRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cacheKeyYearlyTransferReceiver, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlStatistic)
}
