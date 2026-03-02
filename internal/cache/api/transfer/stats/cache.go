package transfer_stats_cache

import "MamangRust/paymentgatewaygrpc/internal/cache"

type TransferStatsCache interface {
	TransferStatsAmountCache
	TransferStatsStatusCache
}

type transferStatsCache struct {
	TransferStatsAmountCache
	TransferStatsStatusCache
}

func NewTransferStatsCache(store *cache.CacheStore) TransferStatsCache {
	return &transferStatsCache{
		TransferStatsAmountCache: NewTransferStatsAmountCache(store),
		TransferStatsStatusCache: NewTransferStatsStatusCache(store),
	}
}
