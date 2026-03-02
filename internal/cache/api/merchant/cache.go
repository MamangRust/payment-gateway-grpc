package merchant_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	merchant_stats_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/merchant/stats"
	merchant_stats_byapikey_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/merchant/statsbyapikey"
	merchant_stats_bymerchant_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/merchant/statsbymerchant"
)

type mencache struct {
	MerchantQueryCache
	MerchantCommandCache
	MerchantTransactionCache
	merchant_stats_cache.MerchantStatsCache
	merchant_stats_byapikey_cache.MerchantStatsByApiKeyCache
	merchant_stats_bymerchant_cache.MerchantStatsByMerchantCache
}

type MerchantMencache interface {
	MerchantQueryCache
	MerchantCommandCache
	MerchantTransactionCache
	merchant_stats_cache.MerchantStatsCache
	merchant_stats_byapikey_cache.MerchantStatsByApiKeyCache
	merchant_stats_bymerchant_cache.MerchantStatsByMerchantCache
}

func NewMerchantMencache(cacheStore *cache.CacheStore) MerchantMencache {

	return &mencache{
		MerchantQueryCache:   NewMerchantQueryCache(cacheStore),
		MerchantCommandCache: NewMerchantCommandCache(cacheStore),

		MerchantTransactionCache:     NewMerchantTransactionCache(cacheStore),
		MerchantStatsCache:           merchant_stats_cache.NewMerchantStatsCache(cacheStore),
		MerchantStatsByApiKeyCache:   merchant_stats_byapikey_cache.NewMerchantStatsByApiKeyCache(cacheStore),
		MerchantStatsByMerchantCache: merchant_stats_bymerchant_cache.NewMerchantStatsByMerchantCache(cacheStore),
	}
}
