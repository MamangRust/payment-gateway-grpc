package transfer_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type transferStatsByCardStatusCache struct {
	store *cache.CacheStore
}

func NewTransferStatsByCardStatusCache(store *cache.CacheStore) TransferStatsByCardStatusCache {
	return &transferStatsByCardStatusCache{store: store}
}

func (t *transferStatsByCardStatusCache) GetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusSuccessCardNumberRow, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthTransferStatusSuccessCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardStatusCache) SetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*db.GetMonthTransferStatusSuccessCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardStatusCache) GetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusSuccessCardNumberRow, bool) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessByCardKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransferStatusSuccessCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardStatusCache) SetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*db.GetYearlyTransferStatusSuccessCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusSuccessByCardKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardStatusCache) GetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusFailedCardNumberRow, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthTransferStatusFailedCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transferStatsByCardStatusCache) SetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*db.GetMonthTransferStatusFailedCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

func (t *transferStatsByCardStatusCache) GetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusFailedCardNumberRow, bool) {
	key := fmt.Sprintf(transferYearTransferStatusFailedByCardKey, req.CardNumber, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransferStatusFailedCardNumberRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transferStatsByCardStatusCache) SetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*db.GetYearlyTransferStatusFailedCardNumberRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusFailedByCardKey, req.CardNumber, req.Year)
	cache.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
