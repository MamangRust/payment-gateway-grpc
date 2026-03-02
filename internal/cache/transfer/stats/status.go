package transfer_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transferStatsStatusCache struct {
	store *cache.CacheStore
}

func NewTransferStatsStatusCache(store *cache.CacheStore) TransferStatsStatusCache {
	return &transferStatsStatusCache{store: store}
}

func (t *transferStatsStatusCache) GetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusSuccessRow, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthTransferStatusSuccessRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsStatusCache) SetCachedMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer, data []*db.GetMonthTransferStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsStatusCache) GetCachedYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusSuccessRow, bool) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTransferStatusSuccessRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsStatusCache) SetCachedYearlyTransferStatusSuccess(ctx context.Context, year int, data []*db.GetYearlyTransferStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsStatusCache) GetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusFailedRow, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthTransferStatusFailedRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsStatusCache) SetCachedMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer, data []*db.GetMonthTransferStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsStatusCache) GetCachedYearlyTransferStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusFailedRow, bool) {
	key := fmt.Sprintf(transferYearTransferStatusFailedKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransferStatusFailedRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsStatusCache) SetCachedYearlyTransferStatusFailed(ctx context.Context, year int, data []*db.GetYearlyTransferStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
