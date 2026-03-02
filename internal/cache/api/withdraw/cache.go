package withdraw_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	withdraw_stats_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/withdraw/stats"
	withdraw_stats_bycard_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/withdraw/statsbycard"
)

type WithdrawMencache interface {
	WithdrawQueryCache
	WithdrawCommandCache
	withdraw_stats_cache.WithdrawStatsCache
	withdraw_stats_bycard_cache.WithdrawStatsByCardCache
}

type withdrawmencache struct {
	WithdrawQueryCache
	WithdrawCommandCache
	withdraw_stats_cache.WithdrawStatsCache
	withdraw_stats_bycard_cache.WithdrawStatsByCardCache
}

func NewWithdrawMencache(cacheStore *cache.CacheStore) WithdrawMencache {
	return &withdrawmencache{
		WithdrawQueryCache:       NewWithdrawQueryCache(cacheStore),
		WithdrawCommandCache:     NewWithdrawCommandCache(cacheStore),
		WithdrawStatsCache:       withdraw_stats_cache.NewWithdrawStatsCache(cacheStore),
		WithdrawStatsByCardCache: withdraw_stats_bycard_cache.NewWithdrawStatsByCardCache(cacheStore),
	}
}
