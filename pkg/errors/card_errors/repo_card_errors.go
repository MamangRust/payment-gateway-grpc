package card_errors

import "errors"

var (
	ErrFindAllCardsFailed         = errors.New("failed to find all cards")
	ErrFindActiveCardsFailed      = errors.New("failed to find active cards")
	ErrFindTrashedCardsFailed     = errors.New("failed to find trashed cards")
	ErrFindCardByIdFailed         = errors.New("failed to find card by ID")
	ErrFindCardByUserIdFailed     = errors.New("failed to find card by user ID")
	ErrFindCardByCardNumberFailed = errors.New("failed to find card by card number")

	ErrGetTotalBalancesFailed          = errors.New("failed to get total balances")
	ErrGetTotalTopAmountFailed         = errors.New("failed to get total topup amount")
	ErrGetTotalWithdrawAmountFailed    = errors.New("failed to get total withdraw amount")
	ErrGetTotalTransactionAmountFailed = errors.New("failed to get total transaction amount")
	ErrGetTotalTransferAmountFailed    = errors.New("failed to get total transfer amount")

	ErrGetTotalBalanceByCardFailed            = errors.New("failed to get total balance by card")
	ErrGetTotalTopupAmountByCardFailed        = errors.New("failed to get total topup amount by card")
	ErrGetTotalWithdrawAmountByCardFailed     = errors.New("failed to get total withdraw amount by card")
	ErrGetTotalTransactionAmountByCardFailed  = errors.New("failed to get total transaction amount by card")
	ErrGetTotalTransferAmountBySenderFailed   = errors.New("failed to get total transfer amount by sender")
	ErrGetTotalTransferAmountByReceiverFailed = errors.New("failed to get total transfer amount by receiver")

	ErrGetMonthlyBalanceFailed = errors.New("failed to get monthly balance")
	ErrGetYearlyBalanceFailed  = errors.New("failed to get yearly balance")

	ErrGetMonthlyTopupAmountFailed = errors.New("failed to get monthly topup amount")
	ErrGetYearlyTopupAmountFailed  = errors.New("failed to get yearly topup amount")

	ErrGetMonthlyWithdrawAmountFailed = errors.New("failed to get monthly withdraw amount")
	ErrGetYearlyWithdrawAmountFailed  = errors.New("failed to get yearly withdraw amount")

	ErrGetMonthlyTransactionAmountFailed = errors.New("failed to get monthly transaction amount")
	ErrGetYearlyTransactionAmountFailed  = errors.New("failed to get yearly transaction amount")

	ErrGetMonthlyTransferAmountSenderFailed = errors.New("failed to get monthly transfer amount by sender")
	ErrGetYearlyTransferAmountSenderFailed  = errors.New("failed to get yearly transfer amount by sender")

	ErrGetMonthlyTransferAmountReceiverFailed = errors.New("failed to get monthly transfer amount by receiver")
	ErrGetYearlyTransferAmountReceiverFailed  = errors.New("failed to get yearly transfer amount by receiver")

	ErrGetMonthlyBalanceByCardFailed = errors.New("failed to get monthly balance by card")
	ErrGetYearlyBalanceByCardFailed  = errors.New("failed to get yearly balance by card")

	ErrGetMonthlyTopupAmountByCardFailed = errors.New("failed to get monthly topup amount by card")
	ErrGetYearlyTopupAmountByCardFailed  = errors.New("failed to get yearly topup amount by card")

	ErrGetMonthlyWithdrawAmountByCardFailed = errors.New("failed to get monthly withdraw amount by card")
	ErrGetYearlyWithdrawAmountByCardFailed  = errors.New("failed to get yearly withdraw amount by card")

	ErrGetMonthlyTransactionAmountByCardFailed = errors.New("failed to get monthly transaction amount by card")
	ErrGetYearlyTransactionAmountByCardFailed  = errors.New("failed to get yearly transaction amount by card")

	ErrGetMonthlyTransferAmountBySenderFailed   = errors.New("failed to get monthly transfer amount by sender")
	ErrGetYearlyTransferAmountBySenderFailed    = errors.New("failed to get yearly transfer amount by sender")
	ErrGetMonthlyTransferAmountByReceiverFailed = errors.New("failed to get monthly transfer amount by receiver")
	ErrGetYearlyTransferAmountByReceiverFailed  = errors.New("failed to get yearly transfer amount by receiver")

	ErrCreateCardFailed = errors.New("failed to create card")
	ErrUpdateCardFailed = errors.New("failed to update card")

	ErrTrashCardFailed           = errors.New("failed to trash card")
	ErrRestoreCardFailed         = errors.New("failed to restore card")
	ErrDeleteCardPermanentFailed = errors.New("failed to delete card permanently")

	ErrRestoreAllCardsFailed         = errors.New("failed to restore all cards")
	ErrDeleteAllCardsPermanentFailed = errors.New("failed to delete all cards permanently")

	ErrInvalidCardRequest = errors.New("invalid card request data")
	ErrInvalidCardId      = errors.New("invalid card ID")
	ErrInvalidUserId      = errors.New("invalid user ID")
	ErrInvalidCardNumber  = errors.New("invalid card number")
	ErrCardAlreadyExists  = errors.New("card already exists")
	ErrCardNotFound       = errors.New("card not found")
)
