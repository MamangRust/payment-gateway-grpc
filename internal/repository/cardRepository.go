package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"MamangRust/paymentgatewaygrpc/pkg/randomvcc"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type cardRepository struct {
	db *db.Queries
}

func NewCardRepository(db *db.Queries) CardRepository {
	return &cardRepository{
		db: db,
	}
}

func (r *cardRepository) FindAllCards(ctx context.Context, req *requests.FindAllCards) ([]*db.GetCardsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCardsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	cards, err := r.db.GetCards(ctx, reqDb)

	if err != nil {
		return nil, card_errors.ErrFindAllCardsFailed
	}

	return cards, nil
}

func (r *cardRepository) FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*db.GetActiveCardsWithCountRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveCardsWithCount(ctx, reqDb)

	if err != nil {
		return nil, card_errors.ErrFindActiveCardsFailed
	}

	return res, nil
}

func (r *cardRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*db.GetTrashedCardsWithCountRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedCardsWithCount(ctx, reqDb)

	if err != nil {
		return nil, card_errors.ErrFindTrashedCardsFailed
	}

	return res, nil
}

func (r *cardRepository) FindById(ctx context.Context, card_id int) (*db.GetCardByIDRow, error) {
	res, err := r.db.GetCardByID(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByIdFailed
	}

	return res, nil
}

func (r *cardRepository) FindCardByUserId(ctx context.Context, user_id int) (*db.GetCardByUserIDRow, error) {
	res, err := r.db.GetCardByUserID(ctx, int32(user_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByUserIdFailed
	}

	return res, nil
}

func (r *cardRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	res, err := r.db.GetCardByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return res, nil
}

func (r *cardRepository) GetTotalBalances(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalBalance(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalBalancesFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTopAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalTopupAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTopAmountFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalWithdrawAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawAmountFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransactionAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalTransactionAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransactionAmountFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransferAmount(ctx context.Context) (*int64, error) {
	res, err := r.db.GetTotalTransferAmount(ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalBalanceByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalBalanceByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalBalanceByCardFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTopupAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTopupAmountByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTopupAmountByCardFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalWithdrawAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmountByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawAmountByCardFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransactionAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransactionAmountByCardNumber(ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransactionAmountByCardFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransferAmountBySender(ctx context.Context, senderCardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransferAmountBySender(ctx, senderCardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountBySenderFailed
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransferAmountByReceiver(ctx context.Context, receiverCardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransferAmountByReceiver(ctx, receiverCardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountByReceiverFailed
	}

	return &res, nil
}

func (r *cardRepository) GetMonthlyBalance(ctx context.Context, year int) ([]*db.GetMonthlyBalancesRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalances(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyBalance(ctx context.Context, year int) ([]*db.GetYearlyBalancesRow, error) {
	res, err := r.db.GetYearlyBalances(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyTopupAmount(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyTopupAmount(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountRow, error) {
	res, err := r.db.GetYearlyTopupAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetYearlyWithdrawAmountRow, error) {
	res, err := r.db.GetYearlyWithdrawAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyTransactionAmount(ctx context.Context, year int) ([]*db.GetMonthlyTransactionAmountRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmount(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyTransactionAmount(ctx context.Context, year int) ([]*db.GetYearlyTransactionAmountRow, error) {
	res, err := r.db.GetYearlyTransactionAmount(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountSenderRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountSender(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountSenderFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountSenderRow, error) {
	res, err := r.db.GetYearlyTransferAmountSender(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountSenderFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountReceiverRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountReceiver(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountReceiverFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountReceiverRow, error) {
	res, err := r.db.GetYearlyTransferAmountReceiver(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountReceiverFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyBalancesByCardNumberRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalancesByCardNumber(ctx, db.GetMonthlyBalancesByCardNumberParams{
		Column1:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceByCardFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyBalancesByCardNumberRow, error) {
	res, err := r.db.GetYearlyBalancesByCardNumber(ctx, db.GetYearlyBalancesByCardNumberParams{
		Column1:    req.Year,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceByCardFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTopupAmountByCardNumberRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountByCardNumber(ctx, db.GetMonthlyTopupAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountByCardFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTopupAmountByCardNumberRow, error) {
	res, err := r.db.GetYearlyTopupAmountByCardNumber(ctx, db.GetYearlyTopupAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountByCardFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyWithdrawAmountByCardNumberRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmountByCardNumber(ctx, db.GetMonthlyWithdrawAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountByCardFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyWithdrawAmountByCardNumberRow, error) {
	res, err := r.db.GetYearlyWithdrawAmountByCardNumber(ctx, db.GetYearlyWithdrawAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountByCardFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransactionAmountByCardNumberRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmountByCardNumber(ctx, db.GetMonthlyTransactionAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountByCardFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransactionAmountByCardNumberRow, error) {
	res, err := r.db.GetYearlyTransactionAmountByCardNumber(ctx, db.GetYearlyTransactionAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountByCardFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountBySenderRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountBySender(ctx, db.GetMonthlyTransferAmountBySenderParams{
		Column2:      yearStart,
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountBySenderFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountBySenderRow, error) {
	res, err := r.db.GetYearlyTransferAmountBySender(ctx, db.GetYearlyTransferAmountBySenderParams{
		Column2:      int32(req.Year),
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountBySenderFailed
	}

	return res, nil
}

func (r *cardRepository) GetMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountByReceiverRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountByReceiver(ctx, db.GetMonthlyTransferAmountByReceiverParams{
		Column2:    yearStart,
		TransferTo: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountByReceiverFailed
	}

	return res, nil
}

func (r *cardRepository) GetYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountByReceiverRow, error) {
	res, err := r.db.GetYearlyTransferAmountByReceiver(ctx, db.GetYearlyTransferAmountByReceiverParams{
		Column2:    int32(req.Year),
		TransferTo: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountByReceiverFailed
	}

	return res, nil
}

func (r *cardRepository) CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*db.CreateCardRow, error) {
	number, err := randomvcc.RandomCardNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate card number: %w", err)
	}

	expireDate := pgtype.Date{
		Time:  request.ExpireDate,
		Valid: true,
	}

	req := db.CreateCardParams{
		UserID:       int32(request.UserID),
		CardNumber:   number,
		CardType:     request.CardType,
		ExpireDate:   expireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.CreateCard(ctx, req)
	if err != nil {
		return nil, card_errors.ErrCreateCardFailed
	}

	return res, nil
}

func (r *cardRepository) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	expireDate := pgtype.Date{
		Time:  request.ExpireDate,
		Valid: true,
	}

	req := db.UpdateCardParams{
		CardID:       int32(request.CardID),
		CardType:     request.CardType,
		ExpireDate:   expireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.UpdateCard(ctx, req)
	if err != nil {
		return nil, card_errors.ErrUpdateCardFailed
	}

	return res, nil
}

func (r *cardRepository) TrashedCard(ctx context.Context, card_id int) (*db.Card, error) {
	res, err := r.db.TrashCard(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrTrashCardFailed
	}

	return res, nil
}

func (r *cardRepository) RestoreCard(ctx context.Context, card_id int) (*db.Card, error) {
	res, err := r.db.RestoreCard(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrRestoreCardFailed
	}

	return res, nil
}

func (r *cardRepository) DeleteCardPermanent(ctx context.Context, card_id int) (bool, error) {
	err := r.db.DeleteCardPermanently(ctx, int32(card_id))

	if err != nil {
		return false, card_errors.ErrDeleteCardPermanentFailed
	}

	return true, nil
}

func (r *cardRepository) RestoreAllCard(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllCards(ctx)

	if err != nil {
		return false, card_errors.ErrRestoreAllCardsFailed
	}

	return true, nil
}

func (r *cardRepository) DeleteAllCardPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentCards(ctx)

	if err != nil {
		return false, card_errors.ErrDeleteAllCardsPermanentFailed
	}

	return true, nil
}
