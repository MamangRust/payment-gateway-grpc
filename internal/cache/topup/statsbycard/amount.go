package topup_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type topupStatsAmountByCardCache struct {
	store *cache.CacheStore
}

func NewTopupStatsAmountByCardCache(store *cache.CacheStore) TopupStatsAmountByCardCache {
	return &topupStatsAmountByCardCache{store: store}
}

func (s *topupStatsAmountByCardCache) GetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupAmountsByCardNumberRow, bool) {
	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetMonthlyTopupAmountsByCardNumberRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsAmountByCardCache) SetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*db.GetMonthlyTopupAmountsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *topupStatsAmountByCardCache) GetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupAmountsByCardNumberRow, bool) {
	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := cache.GetFromCache[[]*db.GetYearlyTopupAmountsByCardNumberRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *topupStatsAmountByCardCache) SetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*db.GetYearlyTopupAmountsByCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
