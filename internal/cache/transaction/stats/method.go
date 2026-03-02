package transaction_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transactionStatsMethodCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsMethodCache(store *cache.CacheStore) TransactionStatsMethodCache {
	return &transactionStatsMethodCache{store: store}
}

func (t *transactionStatsMethodCache) GetMonthlyPaymentMethodsCache(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsRow, bool) {
	key := fmt.Sprintf(monthTransactionMethodCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyPaymentMethodsRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsMethodCache) SetMonthlyPaymentMethodsCache(ctx context.Context, year int, data []*db.GetMonthlyPaymentMethodsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionMethodCacheKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsMethodCache) GetYearlyPaymentMethodsCache(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodsRow, bool) {
	key := fmt.Sprintf(yearTransactionMethodCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyPaymentMethodsRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsMethodCache) SetYearlyPaymentMethodsCache(ctx context.Context, year int, data []*db.GetYearlyPaymentMethodsRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionMethodCacheKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
