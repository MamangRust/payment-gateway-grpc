package transfer_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
)

type TransferQueryCache interface {
	GetCachedTransfersCache(ctx context.Context, req *requests.FindAllTranfers) (*response.ApiResponsePaginationTransfer, bool)
	SetCachedTransfersCache(ctx context.Context, req *requests.FindAllTranfers, data *response.ApiResponsePaginationTransfer)

	GetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTranfers) (*response.ApiResponsePaginationTransferDeleteAt, bool)
	SetCachedTransferActiveCache(ctx context.Context, req *requests.FindAllTranfers, data *response.ApiResponsePaginationTransferDeleteAt)

	GetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTranfers) (*response.ApiResponsePaginationTransferDeleteAt, bool)
	SetCachedTransferTrashedCache(ctx context.Context, req *requests.FindAllTranfers, data *response.ApiResponsePaginationTransferDeleteAt)

	GetCachedTransferCache(ctx context.Context, id int) (*response.ApiResponseTransfer, bool)
	SetCachedTransferCache(ctx context.Context, data *response.ApiResponseTransfer)

	GetCachedTransferByFrom(ctx context.Context, from string) (*response.ApiResponseTransfers, bool)
	SetCachedTransferByFrom(ctx context.Context, from string, data *response.ApiResponseTransfers)

	GetCachedTransferByTo(ctx context.Context, to string) (*response.ApiResponseTransfers, bool)
	SetCachedTransferByTo(ctx context.Context, to string, data *response.ApiResponseTransfers)
}

type TransferCommandCache interface {
	DeleteTransferCache(ctx context.Context, id int)
}
