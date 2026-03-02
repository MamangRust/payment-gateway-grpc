package saldo_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"context"
	"fmt"
)

type saldoCommandCache struct {
	store *cache.CacheStore
}

func NewSaldoCommandCache(store *cache.CacheStore) SaldoCommandCache {
	return &saldoCommandCache{store: store}
}

func (s *saldoCommandCache) DeleteSaldoCache(ctx context.Context, saldo_id int) {
	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	cache.DeleteFromCache(ctx, s.store, key)
}
