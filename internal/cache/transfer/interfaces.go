package transfer_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
)

type TransferQueryCache interface {
	GetCachedTransfersCache(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTransfersRow, *int, bool)
	SetCachedTransfersCache(ctx context.Context, req *requests.FindAllTranfers, data []*db.GetTransfersRow, total *int)

	GetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetActiveTransfersRow, *int, bool)
	SetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTranfers, data []*db.GetActiveTransfersRow, total *int)

	GetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTrashedTransfersRow, *int, bool)
	SetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTranfers, data []*db.GetTrashedTransfersRow, total *int)

	GetCachedTransferCache(ctx context.Context, id int) (*db.GetTransferByIDRow, bool)
	SetCachedTransferCache(ctx context.Context, data *db.GetTransferByIDRow)

	GetCachedTransferByFrom(ctx context.Context, from string) ([]*db.GetTransfersBySourceCardRow, bool)
	SetCachedTransferByFrom(ctx context.Context, from string, data []*db.GetTransfersBySourceCardRow)

	GetCachedTransferByTo(ctx context.Context, to string) ([]*db.GetTransfersByDestinationCardRow, bool)
	SetCachedTransferByTo(ctx context.Context, to string, data []*db.GetTransfersByDestinationCardRow)
}

type TransferCommandCache interface {
	DeleteTransferCache(ctx context.Context, id int)
}
