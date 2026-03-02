package transaction_stats_bycard_cache

import "MamangRust/paymentgatewaygrpc/internal/cache"

type TransactionStatsByCardCache interface {
	TransactionStatsByCardAmountCache
	TransactionStatsByCardStatusCache
	TransactionStatsByCardMethodCache
}

type transactionStatsByCardCache struct {
	TransactionStatsByCardAmountCache
	TransactionStatsByCardStatusCache
	TransactionStatsByCardMethodCache
}

func NewTransactionStatsByCardCache(store *cache.CacheStore) TransactionStatsByCardCache {
	return &transactionStatsByCardCache{
		TransactionStatsByCardAmountCache: NewTransactionStatsByCardAmountCache(store),
		TransactionStatsByCardStatusCache: NewTransactionStatsByCardStatusCache(store),
		TransactionStatsByCardMethodCache: NewTransactionStatsByCardMethodCache(store),
	}
}
