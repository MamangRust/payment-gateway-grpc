package transaction_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transactionStatsByCardAmountCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByCardAmountCache(store *cache.CacheStore) TransactionStatsByCardAmountCache {
	return &transactionStatsByCardAmountCache{store: store}
}

func (t *transactionStatsByCardAmountCache) GetMonthlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyAmountsByCardNumberRow, bool) {
	key := fmt.Sprintf(monthTransactionAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountsByCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByCardAmountCache) SetMonthlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*db.GetMonthlyAmountsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTransactionAmountByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transactionStatsByCardAmountCache) GetYearlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyAmountsByCardNumberRow, bool) {
	key := fmt.Sprintf(yearTransactionAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyAmountsByCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByCardAmountCache) SetYearlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*db.GetYearlyAmountsByCardNumberRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearTransactionAmountByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
