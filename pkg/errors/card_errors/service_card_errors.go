package card_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"
)

var (
	ErrCardNotFoundRes        = response.NewErrorResponse("Card not found", http.StatusNotFound)
	ErrFailedFindAllCards     = response.NewErrorResponse("Failed to fetch Cards", http.StatusInternalServerError)
	ErrFailedFindActiveCards  = response.NewErrorResponse("Failed to fetch active Cards", http.StatusInternalServerError)
	ErrFailedFindTrashedCards = response.NewErrorResponse("Failed to fetch trashed Cards", http.StatusInternalServerError)
	ErrFailedFindById         = response.NewErrorResponse("Failed to find Card by ID", http.StatusInternalServerError)
	ErrFailedFindByUserID     = response.NewErrorResponse("Failed to find Card by User ID", http.StatusInternalServerError)
	ErrFailedFindByCardNumber = response.NewErrorResponse("Failed to find Card by Card Number", http.StatusInternalServerError)

	ErrFailedFindTotalBalances          = response.NewErrorResponse("Failed to Find total balances", http.StatusInternalServerError)
	ErrFailedFindTotalTopAmount         = response.NewErrorResponse("Failed to Find total topup amount", http.StatusInternalServerError)
	ErrFailedFindTotalWithdrawAmount    = response.NewErrorResponse("Failed to Find total withdraw amount", http.StatusInternalServerError)
	ErrFailedFindTotalTransactionAmount = response.NewErrorResponse("Failed to Find total transaction amount", http.StatusInternalServerError)
	ErrFailedFindTotalTransferAmount    = response.NewErrorResponse("Failed to Find total transfer amount", http.StatusInternalServerError)

	ErrFailedFindTotalBalanceByCard            = response.NewErrorResponse("Failed to Find total balance by card", http.StatusInternalServerError)
	ErrFailedFindTotalTopupAmountByCard        = response.NewErrorResponse("Failed to Find total topup amount by card", http.StatusInternalServerError)
	ErrFailedFindTotalWithdrawAmountByCard     = response.NewErrorResponse("Failed to Find total withdraw amount by card", http.StatusInternalServerError)
	ErrFailedFindTotalTransactionAmountByCard  = response.NewErrorResponse("Failed to Find total transaction amount by card", http.StatusInternalServerError)
	ErrFailedFindTotalTransferAmountBySender   = response.NewErrorResponse("Failed to Find total transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindTotalTransferAmountByReceiver = response.NewErrorResponse("Failed to Find total transfer amount by receiver", http.StatusInternalServerError)

	ErrFailedDashboardCard       = response.NewErrorResponse("Failed to get Card dashboard", http.StatusInternalServerError)
	ErrFailedDashboardCardNumber = response.NewErrorResponse("Failed to get Card dashboard by card number", http.StatusInternalServerError)

	ErrFailedFindMonthlyBalance                = response.NewErrorResponse("Failed to get monthly balance", http.StatusInternalServerError)
	ErrFailedFindYearlyBalance                 = response.NewErrorResponse("Failed to get yearly balance", http.StatusInternalServerError)
	ErrFailedFindMonthlyTopupAmount            = response.NewErrorResponse("Failed to get monthly topup amount", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupAmount             = response.NewErrorResponse("Failed to get yearly topup amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyWithdrawAmount         = response.NewErrorResponse("Failed to get monthly withdraw amount", http.StatusInternalServerError)
	ErrFailedFindYearlyWithdrawAmount          = response.NewErrorResponse("Failed to get yearly withdraw amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransactionAmount      = response.NewErrorResponse("Failed to get monthly transaction amount", http.StatusInternalServerError)
	ErrFailedFindYearlyTransactionAmount       = response.NewErrorResponse("Failed to get yearly transaction amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountSender   = response.NewErrorResponse("Failed to get monthly transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountSender    = response.NewErrorResponse("Failed to get yearly transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountReceiver = response.NewErrorResponse("Failed to get monthly transfer amount by receiver", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountReceiver  = response.NewErrorResponse("Failed to get yearly transfer amount by receiver", http.StatusInternalServerError)

	ErrFailedFindMonthlyBalanceByCard            = response.NewErrorResponse("Failed to get monthly balance by card", http.StatusInternalServerError)
	ErrFailedFindYearlyBalanceByCard             = response.NewErrorResponse("Failed to get yearly balance by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyTopupAmountByCard        = response.NewErrorResponse("Failed to get monthly topup amount by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupAmountByCard         = response.NewErrorResponse("Failed to get yearly topup amount by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyWithdrawAmountByCard     = response.NewErrorResponse("Failed to get monthly withdraw amount by card", http.StatusInternalServerError)
	ErrFailedFindYearlyWithdrawAmountByCard      = response.NewErrorResponse("Failed to get yearly withdraw amount by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransactionAmountByCard  = response.NewErrorResponse("Failed to get monthly transaction amount by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTransactionAmountByCard   = response.NewErrorResponse("Failed to get yearly transaction amount by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountBySender   = response.NewErrorResponse("Failed to get monthly transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountBySender    = response.NewErrorResponse("Failed to get yearly transfer amount by sender", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountByReceiver = response.NewErrorResponse("Failed to get monthly transfer amount by receiver", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountByReceiver  = response.NewErrorResponse("Failed to get yearly transfer amount by receiver", http.StatusInternalServerError)

	ErrFailedCreateCard = response.NewErrorResponse("Failed to create Card", http.StatusInternalServerError)
	ErrFailedUpdateCard = response.NewErrorResponse("Failed to update Card", http.StatusInternalServerError)

	ErrFailedTrashCard   = response.NewErrorResponse("Failed to trash Card", http.StatusInternalServerError)
	ErrFailedRestoreCard = response.NewErrorResponse("Failed to restore Card", http.StatusInternalServerError)
	ErrFailedDeleteCard  = response.NewErrorResponse("Failed to delete Card permanently", http.StatusInternalServerError)

	ErrFailedRestoreAllCards = response.NewErrorResponse("Failed to restore all Cards", http.StatusInternalServerError)
	ErrFailedDeleteAllCards  = response.NewErrorResponse("Failed to delete all Cards permanently", http.StatusInternalServerError)
)
