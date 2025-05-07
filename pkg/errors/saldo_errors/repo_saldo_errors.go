package saldo_errors

import "errors"

var (
	ErrFindAllSaldosFailed         = errors.New("failed to find all saldo records")
	ErrFindActiveSaldosFailed      = errors.New("failed to find active saldo records")
	ErrFindTrashedSaldosFailed     = errors.New("failed to find trashed saldo records")
	ErrFindSaldoByIdFailed         = errors.New("failed to find saldo by ID")
	ErrFindSaldoByCardNumberFailed = errors.New("failed to find saldo by card number")

	ErrGetMonthlyTotalSaldoBalanceFailed = errors.New("failed to get monthly total saldo balance")
	ErrGetYearTotalSaldoBalanceFailed    = errors.New("failed to get yearly total saldo balance")
	ErrGetMonthlySaldoBalancesFailed     = errors.New("failed to get monthly saldo balances")
	ErrGetYearlySaldoBalancesFailed      = errors.New("failed to get yearly saldo balances")

	ErrCreateSaldoFailed         = errors.New("failed to create saldo record")
	ErrUpdateSaldoFailed         = errors.New("failed to update saldo record")
	ErrUpdateSaldoBalanceFailed  = errors.New("failed to update saldo balance")
	ErrUpdateSaldoWithdrawFailed = errors.New("failed to update saldo withdrawal")

	ErrTrashSaldoFailed           = errors.New("failed to trash saldo record")
	ErrRestoreSaldoFailed         = errors.New("failed to restore saldo record")
	ErrDeleteSaldoPermanentFailed = errors.New("failed to delete saldo permanently")

	ErrRestoreAllSaldosFailed         = errors.New("failed to restore all saldo records")
	ErrDeleteAllSaldosPermanentFailed = errors.New("failed to delete all saldo records permanently")
)
