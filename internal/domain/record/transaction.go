package record

type TransactionRecord struct {
	ID              int     `json:"id"`
	CardNumber      string  `json:"card_number"`
	TransactionNo   string  `json:"transaction_no"`
	Amount          int     `json:"amount"`
	PaymentMethod   string  `json:"payment_method"`
	MerchantID      int     `json:"merchant_id"`
	TransactionTime string  `json:"transaction_time"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	DeletedAt       *string `json:"deleted_at"`
}

type TransactionRecordMonthStatusSuccess struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type TransactionRecordYearStatusSuccess struct {
	Year         string `json:"year"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type TransactionRecordMonthStatusFailed struct {
	Year        string `json:"year"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type TransactionRecordYearStatusFailed struct {
	Year        string `json:"year"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type TransactionMonthMethod struct {
	Month             string `json:"month"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int    `json:"total_transactions"`
	TotalAmount       int    `json:"total_amount"`
}

type TransactionYearMethod struct {
	Year              string `json:"year"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int    `json:"total_transactions"`
	TotalAmount       int    `json:"total_amount"`
}

type TransactionMonthAmount struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type TransactionYearlyAmount struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}
