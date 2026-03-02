package card_stats_cache

import "MamangRust/paymentgatewaygrpc/internal/cache"

type CardStatsCache interface {
	CardStatsBalanceCache
	CardStatsTopupCache
	CardStatsTransactionCache
	CardStatsTransferCache
	CardStatsWithdrawCache
}

type mencacheStats struct {
	CardStatsBalanceCache
	CardStatsTopupCache
	CardStatsTransactionCache
	CardStatsTransferCache
	CardStatsWithdrawCache
}

func NewMencacheStats(store *cache.CacheStore) CardStatsCache {
	return &mencacheStats{
		CardStatsBalanceCache:     NewCardStatsBalanceCache(store),
		CardStatsTopupCache:       NewCardStatsTopupCache(store),
		CardStatsTransactionCache: NewCardStatsTransactionCache(store),
		CardStatsTransferCache:    NewCardStatsTransferCache(store),
		CardStatsWithdrawCache:    NewCardStatsWithdrawCache(store),
	}
}
