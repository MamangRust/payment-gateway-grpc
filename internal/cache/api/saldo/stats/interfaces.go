package saldo_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
)

type SaldoStatsTotalCache interface {
	GetMonthlyTotalSaldoBalanceCache(ctx context.Context, req *requests.MonthTotalSaldoBalance) (*response.ApiResponseMonthTotalSaldo, bool)
	SetMonthlyTotalSaldoCache(ctx context.Context, req *requests.MonthTotalSaldoBalance, data *response.ApiResponseMonthTotalSaldo)

	GetYearTotalSaldoBalanceCache(ctx context.Context, year int) (*response.ApiResponseYearTotalSaldo, bool)
	SetYearTotalSaldoBalanceCache(ctx context.Context, year int, data *response.ApiResponseYearTotalSaldo)
}

type SaldoStatsBalanceCache interface {
	GetMonthlySaldoBalanceCache(ctx context.Context, year int) (*response.ApiResponseMonthSaldoBalances, bool)
	SetMonthlySaldoBalanceCache(ctx context.Context, year int, data *response.ApiResponseMonthSaldoBalances)

	GetYearlySaldoBalanceCache(ctx context.Context, year int) (*response.ApiResponseYearSaldoBalances, bool)
	SetYearlySaldoBalanceCache(ctx context.Context, year int, data *response.ApiResponseYearSaldoBalances)
}
