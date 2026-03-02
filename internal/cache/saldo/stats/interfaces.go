package saldo_stats_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
)

type SaldoStatsTotalCache interface {
	GetMonthlyTotalSaldoBalanceCache(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*db.GetMonthlyTotalSaldoBalanceRow, bool)
	SetMonthlyTotalSaldoCache(ctx context.Context, req *requests.MonthTotalSaldoBalance, data []*db.GetMonthlyTotalSaldoBalanceRow)

	GetYearTotalSaldoBalanceCache(ctx context.Context, year int) ([]*db.GetYearlyTotalSaldoBalancesRow, bool)
	SetYearTotalSaldoBalanceCache(ctx context.Context, year int, data []*db.GetYearlyTotalSaldoBalancesRow)
}

type SaldoStatsBalanceCache interface {
	GetMonthlySaldoBalanceCache(ctx context.Context, year int) ([]*db.GetMonthlySaldoBalancesRow, bool)
	SetMonthlySaldoBalanceCache(ctx context.Context, year int, data []*db.GetMonthlySaldoBalancesRow)

	GetYearlySaldoBalanceCache(ctx context.Context, year int) ([]*db.GetYearlySaldoBalancesRow, bool)
	SetYearlySaldoBalanceCache(ctx context.Context, year int, data []*db.GetYearlySaldoBalancesRow)
}
