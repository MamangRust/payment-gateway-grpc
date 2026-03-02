package saldo_stats_cache

import "MamangRust/paymentgatewaygrpc/internal/cache"

type SaldoStatsCache interface {
	SaldoStatsBalanceCache
	SaldoStatsTotalCache
}

type saldoStatsCache struct {
	SaldoStatsBalanceCache
	SaldoStatsTotalCache
}

func NewSaldoStatsCache(store *cache.CacheStore) SaldoStatsCache {
	return &saldoStatsCache{
		SaldoStatsBalanceCache: NewSaldoStatsBalanceCache(store),
		SaldoStatsTotalCache:   NewSaldoStatsTotalCache(store),
	}
}
