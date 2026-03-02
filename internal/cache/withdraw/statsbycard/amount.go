package withdraw_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type withdrawStatsByCardAmountCache struct {
	store *cache.CacheStore
}

func NewWithdrawStatsByCardAmountCache(store *cache.CacheStore) WithdrawStatsByCardAmountCache {
	return &withdrawStatsByCardAmountCache{store: store}
}

func (w *withdrawStatsByCardAmountCache) GetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetMonthlyWithdrawsByCardNumberRow, bool) {
	key := fmt.Sprintf(monthWithdrawAmountByCardKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyWithdrawsByCardNumberRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsByCardAmountCache) SetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*db.GetMonthlyWithdrawsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthWithdrawAmountByCardKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsByCardAmountCache) GetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetYearlyWithdrawsByCardNumberRow, bool) {
	key := fmt.Sprintf(yearWithdrawAmountByCardKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyWithdrawsByCardNumberRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsByCardAmountCache) SetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*db.GetYearlyWithdrawsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawAmountByCardKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}
