package merchant_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
)

type MerchantQueryCache interface {
	GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, bool)
	SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants, data []*db.GetMerchantsRow, total *int)

	GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetActiveMerchantsRow, *int, bool)
	SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants, data []*db.GetActiveMerchantsRow, total *int)

	GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetTrashedMerchantsRow, *int, bool)
	SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants, data []*db.GetTrashedMerchantsRow, total *int)

	GetCachedMerchant(ctx context.Context, id int) (*db.GetMerchantByIDRow, bool)
	SetCachedMerchant(ctx context.Context, data *db.GetMerchantByIDRow)

	GetCachedMerchantsByUserId(ctx context.Context, userId int) ([]*db.GetMerchantsByUserIDRow, bool)
	SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*db.GetMerchantsByUserIDRow)

	GetCachedMerchantByApiKey(ctx context.Context, apiKey string) (*db.GetMerchantByApiKeyRow, bool)
	SetCachedMerchantByApiKey(ctx context.Context, apiKey string, data *db.GetMerchantByApiKeyRow)
}

type MerchantTransactionCache interface {
	GetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*db.FindAllTransactionsRow, *int, bool)
	SetCacheAllMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions, data []*db.FindAllTransactionsRow, total *int)

	GetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*db.FindAllTransactionsByMerchantRow, *int, bool)
	SetCacheMerchantTransactions(ctx context.Context, req *requests.FindAllMerchantTransactionsById, data []*db.FindAllTransactionsByMerchantRow, total *int)

	GetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*db.FindAllTransactionsByApikeyRow, *int, bool)
	SetCacheMerchantTransactionApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey, data []*db.FindAllTransactionsByApikeyRow, total *int)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(ctx context.Context, id int)
}
