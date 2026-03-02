package transfer_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	transfer_stats_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/transfer/stats"
	transfer_stats_bycard_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/transfer/statsbycard"
)

type TransferMencache interface {
	TransferQueryCache
	TransferCommandCache
	transfer_stats_cache.TransferStatsCache
	transfer_stats_bycard_cache.TransferStatsByCardCache
}

type transfermencache struct {
	TransferQueryCache
	TransferCommandCache
	transfer_stats_cache.TransferStatsCache
	transfer_stats_bycard_cache.TransferStatsByCardCache
}

func NewTransferMencache(cacheStore *cache.CacheStore) TransferMencache {
	return &transfermencache{
		TransferQueryCache:       NewTransferQueryCache(cacheStore),
		TransferCommandCache:     NewTransferCommandCache(cacheStore),
		TransferStatsCache:       transfer_stats_cache.NewTransferStatsCache(cacheStore),
		TransferStatsByCardCache: transfer_stats_bycard_cache.NewTransferStatsByCardCache(cacheStore),
	}
}
