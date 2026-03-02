package topup_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type topupStatsStatusByCardCache struct {
	store *cache.CacheStore
}

func NewTopupStatsStatusByCardCache(store *cache.CacheStore) TopupStatsStatusByCardCache {
	return &topupStatsStatusByCardCache{store: store}
}

func (s *topupStatsStatusByCardCache) GetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusSuccessCardNumberRow, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthTopupStatusSuccessCardNumberRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsStatusByCardCache) SetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*db.GetMonthTopupStatusSuccessCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsStatusByCardCache) GetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusSuccessCardNumberRow, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupStatusSuccessCardNumberRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsStatusByCardCache) SetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*db.GetYearlyTopupStatusSuccessCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsStatusByCardCache) GetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusFailedCardNumberRow, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthTopupStatusFailedCardNumberRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsStatusByCardCache) SetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*db.GetMonthTopupStatusFailedCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsStatusByCardCache) GetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusFailedCardNumberRow, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupStatusFailedCardNumberRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsStatusByCardCache) SetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*db.GetYearlyTopupStatusFailedCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
