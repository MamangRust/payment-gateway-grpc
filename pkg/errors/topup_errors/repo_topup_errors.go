package topup_errors

import "errors"

var (
	ErrFindAllTopupsFailed          = errors.New("failed to find all topups")
	ErrFindTopupsByActiveFailed     = errors.New("failed to find active topups")
	ErrFindTopupsByTrashedFailed    = errors.New("failed to find trashed topups")
	ErrFindTopupsByCardNumberFailed = errors.New("failed to find topups by card number")
	ErrFindTopupByIdFailed          = errors.New("failed to find topup by ID")

	ErrGetMonthTopupStatusSuccessFailed        = errors.New("failed to get monthly topup status success")
	ErrGetYearlyTopupStatusSuccessFailed       = errors.New("failed to get yearly topup status success")
	ErrGetMonthTopupStatusSuccessByCardFailed  = errors.New("failed to get monthly topup status success by card number")
	ErrGetYearlyTopupStatusSuccessByCardFailed = errors.New("failed to get yearly topup status success by card number")

	ErrGetMonthTopupStatusFailedFailed        = errors.New("failed to get monthly topup status failed")
	ErrGetYearlyTopupStatusFailedFailed       = errors.New("failed to get yearly topup status failed")
	ErrGetMonthTopupStatusFailedByCardFailed  = errors.New("failed to get monthly topup status failed by card number")
	ErrGetYearlyTopupStatusFailedByCardFailed = errors.New("failed to get yearly topup status failed by card number")

	ErrGetMonthlyTopupMethodsFailed       = errors.New("failed to get monthly topup methods")
	ErrGetYearlyTopupMethodsFailed        = errors.New("failed to get yearly topup methods")
	ErrGetMonthlyTopupAmountsFailed       = errors.New("failed to get monthly topup amounts")
	ErrGetYearlyTopupAmountsFailed        = errors.New("failed to get yearly topup amounts")
	ErrGetMonthlyTopupMethodsByCardFailed = errors.New("failed to get monthly topup methods by card number")
	ErrGetYearlyTopupMethodsByCardFailed  = errors.New("failed to get yearly topup methods by card number")
	ErrGetMonthlyTopupAmountsByCardFailed = errors.New("failed to get monthly topup amounts by card number")
	ErrGetYearlyTopupAmountsByCardFailed  = errors.New("failed to get yearly topup amounts by card number")

	ErrCreateTopupFailed       = errors.New("failed to create topup")
	ErrUpdateTopupFailed       = errors.New("failed to update topup")
	ErrUpdateTopupAmountFailed = errors.New("failed to update topup amount")
	ErrUpdateTopupStatusFailed = errors.New("failed to update topup status")

	ErrTrashedTopupFailed            = errors.New("failed to soft-delete (trash) topup")
	ErrRestoreTopupFailed            = errors.New("failed to restore topup")
	ErrDeleteTopupPermanentFailed    = errors.New("failed to permanently delete topup")
	ErrRestoreAllTopupFailed         = errors.New("failed to restore all topups")
	ErrDeleteAllTopupPermanentFailed = errors.New("failed to permanently delete all topups")
)
