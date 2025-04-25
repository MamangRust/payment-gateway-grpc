package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/randomvcc"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type cardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CardRecordMapping
}

func NewCardRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CardRecordMapping) *cardRepository {
	return &cardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cardRepository) FindAllCards(req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCardsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	cards, err := r.db.GetCards(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no cards found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to retrieve cards (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int

	if len(cards) > 0 {
		totalCount = int(cards[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCardsRecord(cards), &totalCount, nil
}

func (r *cardRepository) FindByActive(req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveCardsWithCount(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no active cards found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find active cards (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCardRecordsActive(res), &totalCount, nil

}

func (r *cardRepository) FindByTrashed(req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedCardsWithCount(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no trashed cards found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}
		return nil, nil, fmt.Errorf("failed to find cards transaction (page %d, size %d, search '%s'): %w", req.Page, req.PageSize, req.Search, err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCardRecordsTrashed(res), &totalCount, nil
}

func (r *cardRepository) FindById(card_id int) (*record.CardRecord, error) {
	res, err := r.db.GetCardByID(r.ctx, int32(card_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("card not found with ID: %d", card_id)
		}
		return nil, fmt.Errorf("failed to get card by ID %d: %w", card_id, err)
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) FindCardByUserId(user_id int) (*record.CardRecord, error) {
	res, err := r.db.GetCardByUserID(r.ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("card not found with user_id: %d", user_id)
		}
		return nil, fmt.Errorf("failed to get card by user_id %d: %w", user_id, err)
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) FindCardByCardNumber(card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(r.ctx, card_number)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("card not found with card_number: %s", card_number)
		}
		return nil, fmt.Errorf("failed to get card by card_number %s: %w", card_number, err)
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) GetTotalBalances() (*int64, error) {
	res, err := r.db.GetTotalBalance(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total balance data found")
		}
		return nil, fmt.Errorf("failed to get total balances: %w", err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTopAmount() (*int64, error) {
	res, err := r.db.GetTotalTopupAmount(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total top amount data found")
		}
		return nil, fmt.Errorf("failed to get total top amount: %w", err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalWithdrawAmount() (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmount(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total withdraw amount data found")
		}
		return nil, fmt.Errorf("failed to get total withdraw amount: %w", err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransactionAmount() (*int64, error) {
	res, err := r.db.GetTotalTransactionAmount(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total transaction amount data found")
		}
		return nil, fmt.Errorf("failed to get total transaction amount: %w", err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransferAmount() (*int64, error) {
	res, err := r.db.GetTotalTransferAmount(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total transfer amount data found")
		}
		return nil, fmt.Errorf("failed to get total transfer amount: %w", err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalBalanceByCardNumber(cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalBalanceByCardNumber(r.ctx, cardNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total balance data found for card number %s", cardNumber)
		}
		return nil, fmt.Errorf("failed to get total balance by card number %s: %w", cardNumber, err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTopupAmountByCardNumber(cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTopupAmountByCardNumber(r.ctx, cardNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total topup amount data found for card number %s", cardNumber)
		}
		return nil, fmt.Errorf("failed to get total topup amount by card number %s: %w", cardNumber, err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalWithdrawAmountByCardNumber(cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmountByCardNumber(r.ctx, cardNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total withdraw amount data found for card number %s", cardNumber)
		}
		return nil, fmt.Errorf("failed to get total withdraw amount by card number %s: %w", cardNumber, err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransactionAmountByCardNumber(cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransactionAmountByCardNumber(r.ctx, cardNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total transaction amount data found for card number %s", cardNumber)
		}
		return nil, fmt.Errorf("failed to get total transaction amount by card number %s: %w", cardNumber, err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransferAmountBySender(senderCardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransferAmountBySender(r.ctx, senderCardNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total transfer amount data found for sender card number %s", senderCardNumber)
		}
		return nil, fmt.Errorf("failed to get total transfer amount by sender card number %s: %w", senderCardNumber, err)
	}

	return &res, nil
}

func (r *cardRepository) GetTotalTransferAmountByReceiver(receiverCardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransferAmountByReceiver(r.ctx, receiverCardNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no total transfer amount data found for receiver card number %s", receiverCardNumber)
		}
		return nil, fmt.Errorf("failed to get total transfer amount by receiver card number %s: %w", receiverCardNumber, err)
	}

	return &res, nil
}

func (r *cardRepository) GetMonthlyBalance(year int) ([]*record.CardMonthBalance, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalances(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly balance data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly balance for year %d: %w", year, err)
	}

	return r.mapping.ToMonthlyBalances(res), nil
}

func (r *cardRepository) GetYearlyBalance(year int) ([]*record.CardYearlyBalance, error) {
	res, err := r.db.GetYearlyBalances(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly balance data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly balance for year %d: %w", year, err)
	}

	return r.mapping.ToYearlyBalances(res), nil
}

func (r *cardRepository) GetMonthlyTopupAmount(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmount(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly topup amount data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly topup amount for year %d: %w", year, err)
	}

	return r.mapping.ToMonthlyTopupAmounts(res), nil
}

func (r *cardRepository) GetYearlyTopupAmount(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTopupAmount(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly topup amount data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly topup amount for year %d: %w", year, err)
	}

	return r.mapping.ToYearlyTopupAmounts(res), nil
}

func (r *cardRepository) GetMonthlyWithdrawAmount(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmount(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly withdraw amount data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly withdraw amount for year %d: %w", year, err)
	}

	return r.mapping.ToMonthlyWithdrawAmounts(res), nil
}

func (r *cardRepository) GetYearlyWithdrawAmount(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyWithdrawAmount(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly withdraw amount data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly withdraw amount for year %d: %w", year, err)
	}

	return r.mapping.ToYearlyWithdrawAmounts(res), nil
}

func (r *cardRepository) GetMonthlyTransactionAmount(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmount(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transaction amount data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly transaction amount for year %d: %w", year, err)
	}

	return r.mapping.ToMonthlyTransactionAmounts(res), nil
}

func (r *cardRepository) GetYearlyTransactionAmount(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransactionAmount(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction amount data found for year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly transaction amount for year %d: %w", year, err)
	}

	return r.mapping.ToYearlyTransactionAmounts(res), nil
}

func (r *cardRepository) GetMonthlyTransferAmountSender(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountSender(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transfer amount data found for sender in year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly transfer amount for sender in year %d: %w", year, err)
	}

	return r.mapping.ToMonthlyTransferSenderAmounts(res), nil
}

func (r *cardRepository) GetYearlyTransferAmountSender(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountSender(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer amount data found for sender in year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer amount for sender in year %d: %w", year, err)
	}

	return r.mapping.ToYearlyTransferSenderAmounts(res), nil
}

func (r *cardRepository) GetMonthlyTransferAmountReceiver(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountReceiver(r.ctx, yearStart)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transfer amount data found for receiver in year %d", year)
		}
		return nil, fmt.Errorf("failed to get monthly transfer amount for receiver in year %d: %w", year, err)
	}

	return r.mapping.ToMonthlyTransferReceiverAmounts(res), nil
}

func (r *cardRepository) GetYearlyTransferAmountReceiver(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountReceiver(r.ctx, int32(year))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer amount data found for receiver in year %d", year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer amount for receiver in year %d: %w", year, err)
	}

	return r.mapping.ToYearlyTransferReceiverAmounts(res), nil
}

func (r *cardRepository) GetMonthlyBalancesByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthBalance, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalancesByCardNumber(r.ctx, db.GetMonthlyBalancesByCardNumberParams{
		Column1:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly balance data found for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly balances for card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToMonthlyBalancesCardNumber(res), nil
}

func (r *cardRepository) GetYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearlyBalance, error) {
	res, err := r.db.GetYearlyBalancesByCardNumber(r.ctx, db.GetYearlyBalancesByCardNumberParams{
		Column1:    req.Year,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly balance data found for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly balance for card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToYearlyBalancesCardNumber(res), nil
}

func (r *cardRepository) GetMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountByCardNumber(r.ctx, db.GetMonthlyTopupAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly topup amount data found for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly topup amount for card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToMonthlyTopupAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetYearlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTopupAmountByCardNumber(r.ctx, db.GetYearlyTopupAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly topup amount data found for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly topup amount for card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToYearlyTopupAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmountByCardNumber(r.ctx, db.GetMonthlyWithdrawAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly withdraw amount data found for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly withdraw amount for card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToMonthlyWithdrawAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyWithdrawAmountByCardNumber(r.ctx, db.GetYearlyWithdrawAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly withdraw amount data found for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly withdraw amount for card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToYearlyWithdrawAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmountByCardNumber(r.ctx, db.GetMonthlyTransactionAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transaction amount data found for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly transaction amount for card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToMonthlyTransactionAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransactionAmountByCardNumber(r.ctx, db.GetYearlyTransactionAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transaction amount data found for card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly transaction amount for card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToYearlyTransactionAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetMonthlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountBySender(r.ctx, db.GetMonthlyTransferAmountBySenderParams{
		Column2:      yearStart,
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transfer amount data found for sender card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly transfer amount for sender card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToMonthlyTransferSenderAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountBySender(r.ctx, db.GetYearlyTransferAmountBySenderParams{
		Column2:      int32(req.Year),
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer amount data found for sender card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer amount for sender card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToYearlyTransferSenderAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountByReceiver(r.ctx, db.GetMonthlyTransferAmountByReceiverParams{
		Column2:    yearStart,
		TransferTo: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no monthly transfer amount data found for receiver card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get monthly transfer amount for receiver card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToMonthlyTransferReceiverAmountsByCardNumber(res), nil
}

func (r *cardRepository) GetYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountByReceiver(r.ctx, db.GetYearlyTransferAmountByReceiverParams{
		Column2:    int32(req.Year),
		TransferTo: req.CardNumber,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no yearly transfer amount data found for receiver card number %s in year %d", req.CardNumber, req.Year)
		}
		return nil, fmt.Errorf("failed to get yearly transfer amount for receiver card number %s in year %d: %w", req.CardNumber, req.Year, err)
	}

	return r.mapping.ToYearlyTransferReceiverAmountsByCardNumber(res), nil
}

func (r *cardRepository) CreateCard(request *requests.CreateCardRequest) (*record.CardRecord, error) {
	number, err := randomvcc.RandomCardNumber()

	if err != nil {
		return nil, fmt.Errorf("failed to generate card number: %w", err)
	}

	req := db.CreateCardParams{
		UserID:       int32(request.UserID),
		CardNumber:   number,
		CardType:     request.CardType,
		ExpireDate:   request.ExpireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.CreateCard(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid card data: %w", err)
		}
		return nil, fmt.Errorf("failed to create card: invalid or incomplete card data: %w", err)
	}

	return r.mapping.ToCardRecord(res), nil
}
func (r *cardRepository) UpdateCard(request *requests.UpdateCardRequest) (*record.CardRecord, error) {
	req := db.UpdateCardParams{
		CardID:       int32(request.CardID),
		CardType:     request.CardType,
		ExpireDate:   request.ExpireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.UpdateCard(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("card ID %d not found for update", request.CardID)
		}
		return nil, fmt.Errorf("failed to update card ID %d: card not found or invalid update data", request.CardID)
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) TrashedCard(card_id int) (*record.CardRecord, error) {
	res, err := r.db.TrashCard(r.ctx, int32(card_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("card ID %d not found or already trashed", card_id)
		}
		return nil, fmt.Errorf("failed to trash card ID %d: %w", card_id, err)
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) RestoreCard(card_id int) (*record.CardRecord, error) {
	res, err := r.db.RestoreCard(r.ctx, int32(card_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("card ID %d not found in trash", card_id)
		}
		return nil, fmt.Errorf("failed to restore card ID %d: %w", card_id, err)
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) DeleteCardPermanent(card_id int) (bool, error) {
	err := r.db.DeleteCardPermanently(r.ctx, int32(card_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("card ID %d not found or already deleted", card_id)
		}
		return false, fmt.Errorf("failed to permanently delete card ID %d: %w", card_id, err)
	}

	return true, nil
}

func (r *cardRepository) RestoreAllCard() (bool, error) {
	err := r.db.RestoreAllCards(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed card available to restore")
		}
		return false, fmt.Errorf("failed to restore trashed card: %w", err)
	}

	return true, nil
}

func (r *cardRepository) DeleteAllCardPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentCards(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed card available to delete permanently")
		}
		return false, fmt.Errorf("failed to permanently delete card: %w", err)
	}

	return true, nil
}
