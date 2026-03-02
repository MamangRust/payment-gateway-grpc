package saldo_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
)

type saldoCachedResponseAll struct {
	Data         []*db.GetSaldosRow `json:"data"`
	TotalRecords *int               `json:"total_records"`
}

type saldoCachedResponseActive struct {
	Data         []*db.GetActiveSaldosRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type saldoCachedResponseTrashed struct {
	Data         []*db.GetTrashedSaldosRow `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type saldoQueryCache struct {
	store *cache.CacheStore
}

func NewSaldoQueryCache(store *cache.CacheStore) SaldoQueryCache {
	return &saldoQueryCache{store: store}
}

func (s *saldoQueryCache) GetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetSaldosRow, *int, bool) {
	key := fmt.Sprintf(saldoAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[saldoCachedResponseAll](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) SetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos, data []*db.GetSaldosRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetSaldosRow{}
	}

	key := fmt.Sprintf(saldoAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &saldoCachedResponseAll{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *saldoQueryCache) GetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetActiveSaldosRow, *int, bool) {
	key := fmt.Sprintf(saldoActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[saldoCachedResponseActive](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) SetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos, data []*db.GetActiveSaldosRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetActiveSaldosRow{}
	}

	key := fmt.Sprintf(saldoActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &saldoCachedResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *saldoQueryCache) GetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetTrashedSaldosRow, *int, bool) {
	key := fmt.Sprintf(saldoTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[saldoCachedResponseTrashed](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *saldoQueryCache) SetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos, data []*db.GetTrashedSaldosRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTrashedSaldosRow{}
	}

	key := fmt.Sprintf(saldoTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &saldoCachedResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *saldoQueryCache) GetCachedSaldoById(ctx context.Context, saldo_id int) (*db.GetSaldoByIDRow, bool) {
	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	result, found := cache.GetFromCache[*db.GetSaldoByIDRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *saldoQueryCache) SetCachedSaldoById(ctx context.Context, saldo_id int, data *db.GetSaldoByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoByIdCacheKey, saldo_id)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *saldoQueryCache) GetCachedSaldoByCardNumber(ctx context.Context, card_number string) (*db.Saldo, bool) {
	key := fmt.Sprintf(saldoByCardNumberKey, card_number)
	result, found := cache.GetFromCache[*db.Saldo](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *saldoQueryCache) SetCachedSaldoByCardNumber(ctx context.Context, card_number string, data *db.Saldo) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(saldoByCardNumberKey, card_number)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
