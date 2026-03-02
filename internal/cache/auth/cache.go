package auth_cache

import "MamangRust/paymentgatewaygrpc/internal/cache"

type AuthMencache struct {
	IdentityCache IdentityCache
	LoginCache    LoginCache
}

func NewMencache(cacheStore *cache.CacheStore) *AuthMencache {
	return &AuthMencache{
		IdentityCache: NewidentityCache(cacheStore),
		LoginCache:    NewLoginCache(cacheStore),
	}
}
