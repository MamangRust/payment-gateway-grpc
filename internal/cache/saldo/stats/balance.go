package saldo_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type saldoStatsBalanceCache struct {
	store *cache.CacheStore
}

func NewSaldoStatsBalanceCache(store *cache.CacheStore) SaldoStatsBalanceCache {
	return &saldoStatsBalanceCache{store: store}
}

func (c *saldoStatsBalanceCache) GetMonthlySaldoBalanceCache(ctx context.Context, year int) ([]*db.GetMonthlySaldoBalancesRow, bool) {
	key := fmt.Sprintf(saldoMonthBalanceCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetMonthlySaldoBalancesRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *saldoStatsBalanceCache) SetMonthlySaldoBalanceCache(ctx context.Context, year int, data []*db.GetMonthlySaldoBalancesRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoMonthBalanceCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *saldoStatsBalanceCache) GetYearlySaldoBalanceCache(ctx context.Context, year int) ([]*db.GetYearlySaldoBalancesRow, bool) {
	key := fmt.Sprintf(saldoYearlyBalanceCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlySaldoBalancesRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *saldoStatsBalanceCache) SetYearlySaldoBalanceCache(ctx context.Context, year int, data []*db.GetYearlySaldoBalancesRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoYearlyBalanceCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
