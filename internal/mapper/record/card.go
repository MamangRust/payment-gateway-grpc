package recordmapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
)

type cardRecordMapper struct {
}

func NewCardRecordMapper() *cardRecordMapper {
	return &cardRecordMapper{}
}

func (s *cardRecordMapper) ToCardRecord(card *db.Card) *record.CardRecord {
	var deletedAt *string

	if card.DeletedAt.Valid {
		formatedDeletedAt := card.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.CardRecord{
		ID:           int(card.CardID),
		UserID:       int(card.UserID),
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate.Format("2006-01-02"),
		CVV:          card.Cvv,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    card.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}
}

func (s *cardRecordMapper) ToCardGetAll(card *db.GetCardsRow) *record.CardRecord {
	var deletedAt *string

	if card.DeletedAt.Valid {
		formatedDeletedAt := card.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formatedDeletedAt
	}

	return &record.CardRecord{
		ID:           int(card.CardID),
		UserID:       int(card.UserID),
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate.Format("2006-01-02"),
		CVV:          card.Cvv,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    card.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}
}

func (s *cardRecordMapper) ToCardsRecord(cards []*db.GetCardsRow) []*record.CardRecord {
	var records []*record.CardRecord
	for _, card := range cards {
		records = append(records, s.ToCardGetAll(card))
	}
	return records
}

func (s *cardRecordMapper) ToCardRecordActive(card *db.GetActiveCardsWithCountRow) *record.CardRecord {
	var deletedAt *string

	if card.DeletedAt.Valid {
		formattedDeletedAt := card.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}

	return &record.CardRecord{
		ID:           int(card.CardID),
		UserID:       int(card.UserID),
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate.Format("2006-01-02"),
		CVV:          card.Cvv,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    card.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}
}

func (s *cardRecordMapper) ToCardRecordsActive(cards []*db.GetActiveCardsWithCountRow) []*record.CardRecord {
	var records []*record.CardRecord
	for _, card := range cards {
		records = append(records, s.ToCardRecordActive(card))
	}
	return records
}

func (s *cardRecordMapper) ToCardRecordTrashed(card *db.GetTrashedCardsWithCountRow) *record.CardRecord {
	var deletedAt *string

	if card.DeletedAt.Valid {
		formattedDeletedAt := card.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}

	return &record.CardRecord{
		ID:           int(card.CardID),
		UserID:       int(card.UserID),
		CardNumber:   card.CardNumber,
		CardType:     card.CardType,
		ExpireDate:   card.ExpireDate.Format("2006-01-02"),
		CVV:          card.Cvv,
		CardProvider: card.CardProvider,
		CreatedAt:    card.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    card.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}
}

func (s *cardRecordMapper) ToCardRecordsTrashed(cards []*db.GetTrashedCardsWithCountRow) []*record.CardRecord {
	var records []*record.CardRecord
	for _, card := range cards {
		records = append(records, s.ToCardRecordTrashed(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyBalance(card *db.GetMonthlyBalancesRow) *record.CardMonthBalance {
	return &record.CardMonthBalance{
		Month:        card.Month,
		TotalBalance: int64(card.TotalBalance),
	}
}

func (s *cardRecordMapper) ToMonthlyBalances(cards []*db.GetMonthlyBalancesRow) []*record.CardMonthBalance {
	var records []*record.CardMonthBalance
	for _, card := range cards {
		records = append(records, s.ToMonthlyBalance(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyBalance(card *db.GetYearlyBalancesRow) *record.CardYearlyBalance {
	return &record.CardYearlyBalance{
		Year:         card.Year,
		TotalBalance: card.TotalBalance,
	}
}

func (s *cardRecordMapper) ToYearlyBalances(cards []*db.GetYearlyBalancesRow) []*record.CardYearlyBalance {
	var records []*record.CardYearlyBalance
	for _, card := range cards {
		records = append(records, s.ToYearlyBalance(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyTopupAmount(card *db.GetMonthlyTopupAmountRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalTopupAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyTopupAmounts(cards []*db.GetMonthlyTopupAmountRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyTopupAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyTopupAmount(card *db.GetYearlyTopupAmountRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalTopupAmount,
	}
}

func (s *cardRecordMapper) ToYearlyTopupAmounts(cards []*db.GetYearlyTopupAmountRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyTopupAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyWithdrawAmount(card *db.GetMonthlyWithdrawAmountRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalWithdrawAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyWithdrawAmounts(cards []*db.GetMonthlyWithdrawAmountRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyWithdrawAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyWithdrawAmount(card *db.GetYearlyWithdrawAmountRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalWithdrawAmount,
	}
}

func (s *cardRecordMapper) ToYearlyWithdrawAmounts(cards []*db.GetYearlyWithdrawAmountRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyWithdrawAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyTransactionAmount(card *db.GetMonthlyTransactionAmountRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalTransactionAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyTransactionAmounts(cards []*db.GetMonthlyTransactionAmountRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyTransactionAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyTransactionAmount(card *db.GetYearlyTransactionAmountRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalTransactionAmount,
	}
}

func (s *cardRecordMapper) ToYearlyTransactionAmounts(cards []*db.GetYearlyTransactionAmountRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyTransactionAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyTransferSenderAmount(card *db.GetMonthlyTransferAmountSenderRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalSentAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyTransferSenderAmounts(cards []*db.GetMonthlyTransferAmountSenderRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyTransferSenderAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyTransferSenderAmount(card *db.GetYearlyTransferAmountSenderRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalSentAmount,
	}
}

func (s *cardRecordMapper) ToYearlyTransferSenderAmounts(cards []*db.GetYearlyTransferAmountSenderRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyTransferSenderAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyTransferReceiverAmount(card *db.GetMonthlyTransferAmountReceiverRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalReceivedAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyTransferReceiverAmounts(cards []*db.GetMonthlyTransferAmountReceiverRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyTransferReceiverAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyTransferReceiverAmount(card *db.GetYearlyTransferAmountReceiverRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalReceivedAmount,
	}
}

func (s *cardRecordMapper) ToYearlyTransferReceiverAmounts(cards []*db.GetYearlyTransferAmountReceiverRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyTransferReceiverAmount(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyBalanceCardNumber(card *db.GetMonthlyBalancesByCardNumberRow) *record.CardMonthBalance {
	return &record.CardMonthBalance{
		Month:        card.Month,
		TotalBalance: int64(card.TotalBalance),
	}
}

func (s *cardRecordMapper) ToMonthlyBalancesCardNumber(cards []*db.GetMonthlyBalancesByCardNumberRow) []*record.CardMonthBalance {
	var records []*record.CardMonthBalance
	for _, card := range cards {
		records = append(records, s.ToMonthlyBalanceCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyBalanceCardNumber(card *db.GetYearlyBalancesByCardNumberRow) *record.CardYearlyBalance {
	return &record.CardYearlyBalance{
		Year:         card.Year,
		TotalBalance: card.TotalBalance,
	}
}

func (s *cardRecordMapper) ToYearlyBalancesCardNumber(cards []*db.GetYearlyBalancesByCardNumberRow) []*record.CardYearlyBalance {
	var records []*record.CardYearlyBalance
	for _, card := range cards {
		records = append(records, s.ToYearlyBalanceCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyTopupAmountByCardNumber(card *db.GetMonthlyTopupAmountByCardNumberRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalTopupAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyTopupAmountsByCardNumber(cards []*db.GetMonthlyTopupAmountByCardNumberRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyTopupAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyTopupAmountByCardNumber(card *db.GetYearlyTopupAmountByCardNumberRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalTopupAmount,
	}
}

func (s *cardRecordMapper) ToYearlyTopupAmountsByCardNumber(cards []*db.GetYearlyTopupAmountByCardNumberRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyTopupAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyWithdrawAmountByCardNumber(card *db.GetMonthlyWithdrawAmountByCardNumberRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalWithdrawAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyWithdrawAmountsByCardNumber(cards []*db.GetMonthlyWithdrawAmountByCardNumberRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyWithdrawAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyWithdrawAmountByCardNumber(card *db.GetYearlyWithdrawAmountByCardNumberRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalWithdrawAmount,
	}
}

func (s *cardRecordMapper) ToYearlyWithdrawAmountsByCardNumber(cards []*db.GetYearlyWithdrawAmountByCardNumberRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyWithdrawAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyTransactionAmountByCardNumber(card *db.GetMonthlyTransactionAmountByCardNumberRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalTransactionAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyTransactionAmountsByCardNumber(cards []*db.GetMonthlyTransactionAmountByCardNumberRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyTransactionAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyTransactionAmountByCardNumber(card *db.GetYearlyTransactionAmountByCardNumberRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalTransactionAmount,
	}
}

func (s *cardRecordMapper) ToYearlyTransactionAmountsByCardNumber(cards []*db.GetYearlyTransactionAmountByCardNumberRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyTransactionAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyTransferSenderAmountByCardNumber(card *db.GetMonthlyTransferAmountBySenderRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalSentAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyTransferSenderAmountsByCardNumber(cards []*db.GetMonthlyTransferAmountBySenderRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyTransferSenderAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyTransferSenderAmountByCardNumber(card *db.GetYearlyTransferAmountBySenderRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalSentAmount,
	}
}

func (s *cardRecordMapper) ToYearlyTransferSenderAmountsByCardNumber(cards []*db.GetYearlyTransferAmountBySenderRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyTransferSenderAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToMonthlyTransferReceiverAmountByCardNumber(card *db.GetMonthlyTransferAmountByReceiverRow) *record.CardMonthAmount {
	return &record.CardMonthAmount{
		Month:       card.Month,
		TotalAmount: int64(card.TotalReceivedAmount),
	}
}

func (s *cardRecordMapper) ToMonthlyTransferReceiverAmountsByCardNumber(cards []*db.GetMonthlyTransferAmountByReceiverRow) []*record.CardMonthAmount {
	var records []*record.CardMonthAmount
	for _, card := range cards {
		records = append(records, s.ToMonthlyTransferReceiverAmountByCardNumber(card))
	}
	return records
}

func (s *cardRecordMapper) ToYearlyTransferReceiverAmountByCardNumber(card *db.GetYearlyTransferAmountByReceiverRow) *record.CardYearAmount {
	return &record.CardYearAmount{
		Year:        card.Year,
		TotalAmount: card.TotalReceivedAmount,
	}
}

func (s *cardRecordMapper) ToYearlyTransferReceiverAmountsByCardNumber(cards []*db.GetYearlyTransferAmountByReceiverRow) []*record.CardYearAmount {
	var records []*record.CardYearAmount
	for _, card := range cards {
		records = append(records, s.ToYearlyTransferReceiverAmountByCardNumber(card))
	}
	return records
}
