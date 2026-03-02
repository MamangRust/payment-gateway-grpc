package topup_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type topupStatsStatusCache struct {
	store *cache.CacheStore
}

func NewTopupStatsStatusCache(store *cache.CacheStore) TopupStatsStatusCache {
	return &topupStatsStatusCache{store: store}
}

func (c *topupStatsStatusCache) GetMonthTopupStatusSuccessCache(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusSuccessRow, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthTopupStatusSuccessRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatsStatusCache) SetMonthTopupStatusSuccessCache(ctx context.Context, req *requests.MonthTopupStatus, data []*db.GetMonthTopupStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *topupStatsStatusCache) GetYearlyTopupStatusSuccessCache(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusSuccessRow, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupStatusSuccessRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatsStatusCache) SetYearlyTopupStatusSuccessCache(ctx context.Context, year int, data []*db.GetYearlyTopupStatusSuccessRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *topupStatsStatusCache) GetMonthTopupStatusFailedCache(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusFailedRow, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthTopupStatusFailedRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *topupStatsStatusCache) SetMonthTopupStatusFailedCache(ctx context.Context, req *requests.MonthTopupStatus, data []*db.GetMonthTopupStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}

func (c *topupStatsStatusCache) GetYearlyTopupStatusFailedCache(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusFailedRow, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupStatusFailedRow](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (c *topupStatsStatusCache) SetYearlyTopupStatusFailedCache(ctx context.Context, year int, data []*db.GetYearlyTopupStatusFailedRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)
	cache.SetToCache(ctx, c.store, key, &data, ttlDefault)
}
