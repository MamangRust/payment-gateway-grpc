package topup_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	topup_stats_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/topup/stats"
	topup_stats_bycard_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/topup/statsbycard"
)

type TopupMencach interface {
	TopupQueryCache
	TopupCommandCache
	topup_stats_cache.TopupStatsCache
	topup_stats_bycard_cache.TopupStatsByCardCache
}

type mencache struct {
	TopupQueryCache
	TopupCommandCache
	topup_stats_cache.TopupStatsCache
	topup_stats_bycard_cache.TopupStatsByCardCache
}

func NewTopupMencache(cacheStore *cache.CacheStore) TopupMencach {

	return &mencache{
		TopupQueryCache:       NewTopupQueryCache(cacheStore),
		TopupCommandCache:     NewTopupCommandCache(cacheStore),
		TopupStatsCache:       topup_stats_cache.NewTopupStatsCache(cacheStore),
		TopupStatsByCardCache: topup_stats_bycard_cache.NewTopupStatsByCardCache(cacheStore),
	}
}
