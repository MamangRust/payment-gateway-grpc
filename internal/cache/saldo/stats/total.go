package saldo_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type saldoStatsTotalCache struct {
	store *cache.CacheStore
}

func NewSaldoStatsTotalCache(store *cache.CacheStore) SaldoStatsTotalCache {
	return &saldoStatsTotalCache{store: store}
}

func (c *saldoStatsTotalCache) GetMonthlyTotalSaldoBalanceCache(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*db.GetMonthlyTotalSaldoBalanceRow, bool) {
	key := fmt.Sprintf(saldoMonthTotalBalanceCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalSaldoBalanceRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *saldoStatsTotalCache) SetMonthlyTotalSaldoCache(ctx context.Context, req *requests.MonthTotalSaldoBalance, data []*db.GetMonthlyTotalSaldoBalanceRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoMonthTotalBalanceCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *saldoStatsTotalCache) GetYearTotalSaldoBalanceCache(ctx context.Context, year int) ([]*db.GetYearlyTotalSaldoBalancesRow, bool) {
	key := fmt.Sprintf(saldoYearTotalBalanceCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTotalSaldoBalancesRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *saldoStatsTotalCache) SetYearTotalSaldoBalanceCache(ctx context.Context, year int, data []*db.GetYearlyTotalSaldoBalancesRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoYearTotalBalanceCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
