package withdraw_stats_bycard_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
)

type WithdrawStatsByCardStatusCache interface {
	GetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusSuccessCardNumberRow, bool)
	SetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*db.GetMonthWithdrawStatusSuccessCardNumberRow)

	GetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusSuccessCardNumberRow, bool)
	SetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*db.GetYearlyWithdrawStatusSuccessCardNumberRow)

	GetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusFailedCardNumberRow, bool)
	SetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*db.GetMonthWithdrawStatusFailedCardNumberRow)

	GetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusFailedCardNumberRow, bool)
	SetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*db.GetYearlyWithdrawStatusFailedCardNumberRow)
}

type WithdrawStatsByCardAmountCache interface {
	GetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetMonthlyWithdrawsByCardNumberRow, bool)
	SetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*db.GetMonthlyWithdrawsByCardNumberRow)

	GetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetYearlyWithdrawsByCardNumberRow, bool)
	SetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*db.GetYearlyWithdrawsByCardNumberRow)
}
