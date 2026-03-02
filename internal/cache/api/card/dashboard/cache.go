package card_dashboard_cache

import "MamangRust/paymentgatewaygrpc/internal/cache"

type CardDashboardCache interface {
	CardDashboardTotalCache
	CardDashboardByCardNumberCache
}

type cardDashboardCaches struct {
	CardDashboardTotalCache
	CardDashboardByCardNumberCache
}

func NewMencacheDashboard(store *cache.CacheStore) CardDashboardCache {
	return &cardDashboardCaches{
		CardDashboardTotalCache:        NewCardDashboardCache(store),
		CardDashboardByCardNumberCache: NewCardDashboardByCardNumberCache(store),
	}
}
