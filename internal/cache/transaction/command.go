package transaction_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"context"
	"fmt"
)

type transactionCommandCache struct {
	store *cache.CacheStore
}

func NewTransactionCommandCache(store *cache.CacheStore) TransactionCommandCache {
	return &transactionCommandCache{store: store}
}

func (t *transactionCommandCache) DeleteTransactionCache(ctx context.Context, id int) {
	cache.DeleteFromCache(ctx, t.store, fmt.Sprintf(transactionByIdCacheKey, id))
}
