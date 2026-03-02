package saldo_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
)

type SaldoQueryCache interface {
	GetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetSaldosRow, *int, bool)
	SetCachedSaldos(ctx context.Context, req *requests.FindAllSaldos, data []*db.GetSaldosRow, totalRecords *int)

	GetCachedSaldoById(ctx context.Context, saldo_id int) (*db.GetSaldoByIDRow, bool)
	SetCachedSaldoById(ctx context.Context, saldo_id int, data *db.GetSaldoByIDRow)

	GetCachedSaldoByCardNumber(ctx context.Context, card_number string) (*db.Saldo, bool)
	SetCachedSaldoByCardNumber(ctx context.Context, card_number string, data *db.Saldo)

	GetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetActiveSaldosRow, *int, bool)
	SetCachedSaldoByActive(ctx context.Context, req *requests.FindAllSaldos, data []*db.GetActiveSaldosRow, totalRecords *int)

	GetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetTrashedSaldosRow, *int, bool)
	SetCachedSaldoByTrashed(ctx context.Context, req *requests.FindAllSaldos, data []*db.GetTrashedSaldosRow, totalRecords *int)
}

type SaldoCommandCache interface {
	DeleteSaldoCache(ctx context.Context, saldo_id int)
}
