package card_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"net/http"
)

var (
	ErrCardNotFoundRes        = errors.NewErrorResponse("Card not found", http.StatusNotFound)
	ErrFailedFindAllCards     = errors.NewErrorResponse("Failed to fetch Cards", http.StatusInternalServerError)
	ErrFailedFindActiveCards  = errors.NewErrorResponse("Failed to fetch active Cards", http.StatusInternalServerError)
	ErrFailedFindTrashedCards = errors.NewErrorResponse("Failed to fetch trashed Cards", http.StatusInternalServerError)
	ErrFailedFindById         = errors.NewErrorResponse("Failed to find Card by ID", http.StatusInternalServerError)
	ErrFailedFindByUserID     = errors.NewErrorResponse("Failed to find Card by User ID", http.StatusInternalServerError)
	ErrFailedFindByCardNumber = errors.NewErrorResponse("Failed to find Card by Card Number", http.StatusInternalServerError)

	ErrFailedFindTotalBalances          = errors.NewErrorResponse("Failed to Find total balances", http.StatusInternalServerError)
	ErrFailedFindTotalTopAmount         = errors.NewErrorResponse("Failed to Find total topup amount", http.StatusInternalServerError)
	ErrFailedFindTotalWithdrawAmount    = errors.NewErrorResponse("Failed to Find total withdraw amount", http.StatusInternalServerError)
	ErrFailedFindTotalTransactionAmount = errors.NewErrorResponse("Failed to Find total transaction amount", http.StatusInternalServerError)
	ErrFailedFindTotalTransferAmount    = errors.NewErrorResponse("Failed to Find total transfer amount", http.StatusInternalServerError)

	ErrFailedFindTotalBalanceByCard            = errors.NewErrorResponse("Failed to Find total balance by card", http.StatusInternalServerError)
	ErrFailedFindTotalTopupAmountByCard        = errors.NewErrorResponse("Failed to Find total topup amount by card", http.StatusInternalServerError)
	ErrFailedFindTotalWithdrawAmountByCard     = errors.NewErrorResponse("Failed to Find total withdraw amount by card", http.StatusInternalServerError)
	ErrFailedFindTotalTransactionAmountByCard  = errors.NewErrorResponse("Failed to Find total transaction amount by card", http.StatusInternalServerError)
	ErrFailedFindTotalTransferAmountBySender   = errors.NewErrorResponse("Failed to Find total transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindTotalTransferAmountByReceiver = errors.NewErrorResponse("Failed to Find total transfer amount by receiver", http.StatusInternalServerError)

	ErrFailedDashboardCard       = errors.NewErrorResponse("Failed to get Card dashboard", http.StatusInternalServerError)
	ErrFailedDashboardCardNumber = errors.NewErrorResponse("Failed to get Card dashboard by card number", http.StatusInternalServerError)

	ErrFailedFindMonthlyBalance                = errors.NewErrorResponse("Failed to get monthly balance", http.StatusInternalServerError)
	ErrFailedFindYearlyBalance                 = errors.NewErrorResponse("Failed to get yearly balance", http.StatusInternalServerError)
	ErrFailedFindMonthlyTopupAmount            = errors.NewErrorResponse("Failed to get monthly topup amount", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupAmount             = errors.NewErrorResponse("Failed to get yearly topup amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyWithdrawAmount         = errors.NewErrorResponse("Failed to get monthly withdraw amount", http.StatusInternalServerError)
	ErrFailedFindYearlyWithdrawAmount          = errors.NewErrorResponse("Failed to get yearly withdraw amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransactionAmount      = errors.NewErrorResponse("Failed to get monthly transaction amount", http.StatusInternalServerError)
	ErrFailedFindYearlyTransactionAmount       = errors.NewErrorResponse("Failed to get yearly transaction amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountSender   = errors.NewErrorResponse("Failed to get monthly transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountSender    = errors.NewErrorResponse("Failed to get yearly transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountReceiver = errors.NewErrorResponse("Failed to get monthly transfer amount by receiver", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountReceiver  = errors.NewErrorResponse("Failed to get yearly transfer amount by receiver", http.StatusInternalServerError)

	ErrFailedFindMonthlyBalanceByCard            = errors.NewErrorResponse("Failed to get monthly balance by card", http.StatusInternalServerError)
	ErrFailedFindYearlyBalanceByCard             = errors.NewErrorResponse("Failed to get yearly balance by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyTopupAmountByCard        = errors.NewErrorResponse("Failed to get monthly topup amount by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupAmountByCard         = errors.NewErrorResponse("Failed to get yearly topup amount by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyWithdrawAmountByCard     = errors.NewErrorResponse("Failed to get monthly withdraw amount by card", http.StatusInternalServerError)
	ErrFailedFindYearlyWithdrawAmountByCard      = errors.NewErrorResponse("Failed to get yearly withdraw amount by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransactionAmountByCard  = errors.NewErrorResponse("Failed to get monthly transaction amount by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTransactionAmountByCard   = errors.NewErrorResponse("Failed to get yearly transaction amount by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountBySender   = errors.NewErrorResponse("Failed to get monthly transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountBySender    = errors.NewErrorResponse("Failed to get yearly transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountByReceiver = errors.NewErrorResponse("Failed to get monthly transfer amount by receiver", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountByReceiver  = errors.NewErrorResponse("Failed to get yearly transfer amount by receiver", http.StatusInternalServerError)

	ErrFailedCreateCard = errors.NewErrorResponse("Failed to create Card", http.StatusInternalServerError)
	ErrFailedUpdateCard = errors.NewErrorResponse("Failed to update Card", http.StatusInternalServerError)

	ErrFailedTrashCard   = errors.NewErrorResponse("Failed to trash Card", http.StatusInternalServerError)
	ErrFailedRestoreCard = errors.NewErrorResponse("Failed to restore Card", http.StatusInternalServerError)
	ErrFailedDeleteCard  = errors.NewErrorResponse("Failed to delete Card permanently", http.StatusInternalServerError)

	ErrFailedRestoreAllCards = errors.NewErrorResponse("Failed to restore all Cards", http.StatusInternalServerError)
	ErrFailedDeleteAllCards  = errors.NewErrorResponse("Failed to delete all Cards permanently", http.StatusInternalServerError)
)
