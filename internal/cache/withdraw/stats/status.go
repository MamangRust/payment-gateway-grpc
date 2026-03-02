package withdraw_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type withdrawStatsStatusCache struct {
	store *cache.CacheStore
}

func NewWithdrawStatsStatusCache(store *cache.CacheStore) WithdrawStatsStatusCache {
	return &withdrawStatsStatusCache{store: store}
}

func (w *withdrawStatsStatusCache) GetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusSuccessRow, bool) {
	key := fmt.Sprintf(montWithdrawStatusSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthWithdrawStatusSuccessRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsStatusCache) SetCachedMonthWithdrawStatusSuccessCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*db.GetMonthWithdrawStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawStatusSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsStatusCache) GetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusSuccessRow, bool) {
	key := fmt.Sprintf(yearWithdrawStatusSuccessKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyWithdrawStatusSuccessRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsStatusCache) SetCachedYearlyWithdrawStatusSuccessCache(ctx context.Context, year int, data []*db.GetYearlyWithdrawStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusSuccessKey, year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsStatusCache) GetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusFailedRow, bool) {
	key := fmt.Sprintf(montWithdrawStatusFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthWithdrawStatusFailedRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsStatusCache) SetCachedMonthWithdrawStatusFailedCache(ctx context.Context, req *requests.MonthStatusWithdraw, data []*db.GetMonthWithdrawStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(montWithdrawStatusFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}

func (w *withdrawStatsStatusCache) GetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusFailedRow, bool) {
	key := fmt.Sprintf(yearWithdrawStatusFailedKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyWithdrawStatusFailedRow](ctx, w.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (w *withdrawStatsStatusCache) SetCachedYearlyWithdrawStatusFailedCache(ctx context.Context, year int, data []*db.GetYearlyWithdrawStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearWithdrawStatusFailedKey, year)
	cache.SetToCache(ctx, w.store, key, &data, ttlDefault)
}
