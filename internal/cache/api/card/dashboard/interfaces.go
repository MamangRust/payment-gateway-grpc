package card_dashboard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
)

type CardDashboardTotalCache interface {
	GetDashboardCardCache(ctx context.Context) (*response.ApiResponseDashboardCard, bool)
	SetDashboardCardCache(ctx context.Context, data *response.ApiResponseDashboardCard)
	DeleteDashboardCardCache(ctx context.Context)
}

type CardDashboardByCardNumberCache interface {
	SetDashboardCardCardNumberCache(ctx context.Context, cardNumber string, data *response.ApiResponseDashboardCardNumber)
	GetDashboardCardCardNumberCache(ctx context.Context, cardNumber string) (*response.ApiResponseDashboardCardNumber, bool)
	DeleteDashboardCardCardNumberCache(ctx context.Context, cardNumber string)
}
