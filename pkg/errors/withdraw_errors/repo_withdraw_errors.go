package withdraw_errors

import "errors"

var (
	ErrFindAllWithdrawsFailed          = errors.New("failed to find all withdraws")
	ErrFindActiveWithdrawsFailed       = errors.New("failed to find active withdraws")
	ErrFindTrashedWithdrawsFailed      = errors.New("failed to find trashed withdraws")
	ErrFindWithdrawsByCardNumberFailed = errors.New("failed to find withdraws by card number")
	ErrFindWithdrawByIdFailed          = errors.New("failed to find withdraw by ID")

	ErrGetMonthWithdrawStatusSuccessFailed        = errors.New("failed to get monthly withdraw status success")
	ErrGetYearlyWithdrawStatusSuccessFailed       = errors.New("failed to get yearly withdraw status success")
	ErrGetMonthWithdrawStatusSuccessByCardFailed  = errors.New("failed to get monthly withdraw status success by card number")
	ErrGetYearlyWithdrawStatusSuccessByCardFailed = errors.New("failed to get yearly withdraw status success by card number")

	ErrGetMonthWithdrawStatusFailedFailed        = errors.New("failed to get monthly withdraw status failed")
	ErrGetYearlyWithdrawStatusFailedFailed       = errors.New("failed to get yearly withdraw status failed")
	ErrGetMonthWithdrawStatusFailedByCardFailed  = errors.New("failed to get monthly withdraw status failed by card number")
	ErrGetYearlyWithdrawStatusFailedByCardFailed = errors.New("failed to get yearly withdraw status failed by card number")

	ErrGetMonthlyWithdrawsFailed       = errors.New("failed to get monthly withdraw amounts")
	ErrGetYearlyWithdrawsFailed        = errors.New("failed to get yearly withdraw amounts")
	ErrGetMonthlyWithdrawsByCardFailed = errors.New("failed to get monthly withdraw amounts by card number")
	ErrGetYearlyWithdrawsByCardFailed  = errors.New("failed to get yearly withdraw amounts by card number")

	ErrCreateWithdrawFailed       = errors.New("failed to create withdraw")
	ErrUpdateWithdrawFailed       = errors.New("failed to update withdraw")
	ErrUpdateWithdrawStatusFailed = errors.New("failed to update withdraw status")

	ErrTrashedWithdrawFailed             = errors.New("failed to soft-delete (trash) withdraw")
	ErrRestoreWithdrawFailed             = errors.New("failed to restore withdraw")
	ErrDeleteWithdrawPermanentFailed     = errors.New("failed to permanently delete withdraw")
	ErrRestoreAllWithdrawsFailed         = errors.New("failed to restore all withdraws")
	ErrDeleteAllWithdrawsPermanentFailed = errors.New("failed to permanently delete all withdraws")
)
