package transaction_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transactionStatsStatusCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsStatusCache(store *cache.CacheStore) TransactionStatsStatusCache {
	return &transactionStatsStatusCache{store: store}
}

func (t *transactionStatsStatusCache) GetMonthTransactionStatusSuccessCache(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusSuccessRow, bool) {
	key := fmt.Sprintf(monthTransactionStatusSuccessCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthTransactionStatusSuccessRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsStatusCache) SetMonthTransactionStatusSuccessCache(ctx context.Context, req *requests.MonthStatusTransaction, data []*db.GetMonthTransactionStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionStatusSuccessCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsStatusCache) GetYearTransactionStatusSuccessCache(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusSuccessRow, bool) {
	key := fmt.Sprintf(yearTransactionStatusSuccessCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionStatusSuccessRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsStatusCache) SetYearTransactionStatusSuccessCache(ctx context.Context, year int, data []*db.GetYearlyTransactionStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionStatusSuccessCacheKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsStatusCache) GetMonthTransactionStatusFailedCache(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusFailedRow, bool) {
	key := fmt.Sprintf(monthTransactionStatusFailedCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthTransactionStatusFailedRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsStatusCache) SetMonthTransactionStatusFailedCache(ctx context.Context, req *requests.MonthStatusTransaction, data []*db.GetMonthTransactionStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionStatusFailedCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsStatusCache) GetYearTransactionStatusFailedCache(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusFailedRow, bool) {
	key := fmt.Sprintf(yearTransactionStatusFailedCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionStatusFailedRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsStatusCache) SetYearTransactionStatusFailedCache(ctx context.Context, year int, data []*db.GetYearlyTransactionStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionStatusFailedCacheKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
