package topup_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type topupStatsMethodByCardCache struct {
	store *cache.CacheStore
}

func NewTopupStatsMethodByCardCache(store *cache.CacheStore) TopupStatsMethodByCardCache {
	return &topupStatsMethodByCardCache{store: store}
}

func (s *topupStatsMethodByCardCache) GetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupMethodsByCardNumberRow, bool) {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTopupMethodsByCardNumberRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsMethodByCardCache) SetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*db.GetMonthlyTopupMethodsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsMethodByCardCache) GetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupMethodsByCardNumberRow, bool) {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupMethodsByCardNumberRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsMethodByCardCache) SetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*db.GetYearlyTopupMethodsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
