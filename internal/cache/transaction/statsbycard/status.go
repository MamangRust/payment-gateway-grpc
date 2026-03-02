package transaction_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transactionStatsByCardStatusCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByCardStatusCache(store *cache.CacheStore) TransactionStatsByCardStatusCache {
	return &transactionStatsByCardStatusCache{store: store}
}

func (t *transactionStatsByCardStatusCache) GetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusSuccessCardNumberRow, bool) {
	key := fmt.Sprintf(monthTransactionStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthTransactionStatusSuccessCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByCardStatusCache) SetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*db.GetMonthTransactionStatusSuccessCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardStatusCache) GetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusSuccessCardNumberRow, bool) {
	key := fmt.Sprintf(yearTransactionStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionStatusSuccessCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByCardStatusCache) SetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*db.GetYearlyTransactionStatusSuccessCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardStatusCache) GetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusFailedCardNumberRow, bool) {
	key := fmt.Sprintf(monthTransactionStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthTransactionStatusFailedCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByCardStatusCache) SetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*db.GetMonthTransactionStatusFailedCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardStatusCache) GetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusFailedCardNumberRow, bool) {
	key := fmt.Sprintf(yearTransactionStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionStatusFailedCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByCardStatusCache) SetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*db.GetYearlyTransactionStatusFailedCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
