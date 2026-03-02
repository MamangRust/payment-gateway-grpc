package transaction_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transactionStatsByCardMethodCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByCardMethodCache(store *cache.CacheStore) TransactionStatsByCardMethodCache {
	return &transactionStatsByCardMethodCache{store: store}
}

func (t *transactionStatsByCardMethodCache) GetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyPaymentMethodsByCardNumberRow, bool) {
	key := fmt.Sprintf(monthTransactionMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyPaymentMethodsByCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByCardMethodCache) SetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*db.GetMonthlyPaymentMethodsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionMethodByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardMethodCache) GetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyPaymentMethodsByCardNumberRow, bool) {
	key := fmt.Sprintf(yearTransactionMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyPaymentMethodsByCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByCardMethodCache) SetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*db.GetYearlyPaymentMethodsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTransactionMethodByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
