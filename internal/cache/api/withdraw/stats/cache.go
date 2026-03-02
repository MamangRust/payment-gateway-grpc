package withdraw_stats_cache

import "MamangRust/paymentgatewaygrpc/internal/cache"

type WithdrawStatsCache interface {
	WithdrawStatsAmountCache
	WithdrawStatsStatusCache
}

type withdrawStatsCache struct {
	WithdrawStatsAmountCache
	WithdrawStatsStatusCache
}

func NewWithdrawStatsCache(store *cache.CacheStore) WithdrawStatsCache {
	return &withdrawStatsCache{
		WithdrawStatsAmountCache: NewWithdrawStatsAmountCache(store),
		WithdrawStatsStatusCache: NewWithdrawStatsStatusCache(store),
	}
}
