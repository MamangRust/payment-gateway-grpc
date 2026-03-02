package transaction_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	transaction_stats_cache "MamangRust/paymentgatewaygrpc/internal/cache/transaction/stats"
	transaction_stats_bycard_cache "MamangRust/paymentgatewaygrpc/internal/cache/transaction/statsbycard"
)

type TransactionMencache interface {
	TransactionQueryCache
	TransactionCommandCache
	transaction_stats_cache.TransactionStatsCache
	transaction_stats_bycard_cache.TransactionStatsByCardCache
}

type transactionmencache struct {
	TransactionQueryCache
	TransactionCommandCache
	transaction_stats_cache.TransactionStatsCache
	transaction_stats_bycard_cache.TransactionStatsByCardCache
}

func NewTransactionMencache(cacheStore *cache.CacheStore) TransactionMencache {
	return &transactionmencache{
		TransactionQueryCache:       NewTransactionQueryCache(cacheStore),
		TransactionCommandCache:     NewTransactionCommandCache(cacheStore),
		TransactionStatsCache:       transaction_stats_cache.NewTransactionStatsCache(cacheStore),
		TransactionStatsByCardCache: transaction_stats_bycard_cache.NewTransactionStatsByCardCache(cacheStore),
	}
}
