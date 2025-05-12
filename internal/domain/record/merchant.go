package record

import (
	"time"
)

type MerchantRecord struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ApiKey    string  `json:"api_key"`
	UserID    int     `json:"user_id"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type MerchantTransactionsRecord struct {
	TransactionID   int32     `json:"transaction_id"`
	CardNumber      string    `json:"card_number"`
	Amount          int32     `json:"amount"`
	PaymentMethod   string    `json:"payment_method"`
	MerchantID      int32     `json:"merchant_id"`
	MerchantName    string    `json:"merchant_name"`
	TransactionTime time.Time `json:"transaction_time"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	DeletedAt       *string   `json:"deleted_at"`
}

type MerchantYearlyPaymentMethod struct {
	Year          string `json:"year"`
	PaymentMethod string `json:"payment_method"`
	TotalAmount   int    `json:"total_amount"`
}

type MerchantMonthlyPaymentMethod struct {
	Month         string `json:"month"`
	PaymentMethod string `json:"payment_method"`
	TotalAmount   int    `json:"total_amount"`
}

type MerchantMonthlyAmount struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type MerchantYearlyAmount struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}

type MerchantMonthlyTotalAmount struct {
	Year        string `json:"year"`
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type MerchantYearlyTotalAmount struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}
