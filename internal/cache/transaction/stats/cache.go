package transaction_stats_cache

import "MamangRust/paymentgatewaygrpc/internal/cache"

type TransactionStatsCache interface {
	TransactionStatsAmountCache
	TransactionStatsMethodCache
	TransactionStatsStatusCache
}

type transactonStatsCache struct {
	TransactionStatsAmountCache
	TransactionStatsMethodCache
	TransactionStatsStatusCache
}

func NewTransactionStatsCache(store *cache.CacheStore) TransactionStatsCache {
	return &transactonStatsCache{
		TransactionStatsAmountCache: NewTransactionStatsAmountCache(store),
		TransactionStatsMethodCache: NewTransactionStatsMethodCache(store),
		TransactionStatsStatusCache: NewTransactionStatsStatusCache(store),
	}
}
