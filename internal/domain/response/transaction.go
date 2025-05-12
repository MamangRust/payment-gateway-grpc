package response

type TransactionResponse struct {
	ID              int    `json:"id"`
	TransactionNo   string `json:"transaction_no"`
	CardNumber      string `json:"card_number"`
	Amount          int    `json:"amount"`
	PaymentMethod   string `json:"payment_method"`
	MerchantID      int    `json:"merchant_id"`
	TransactionTime string `json:"transaction_time"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type TransactionResponseDeleteAt struct {
	ID              int     `json:"id"`
	TransactionNo   string  `json:"transaction_no"`
	CardNumber      string  `json:"card_number"`
	Amount          int     `json:"amount"`
	PaymentMethod   string  `json:"payment_method"`
	MerchantID      int     `json:"merchant_id"`
	TransactionTime string  `json:"transaction_time"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	DeletedAt       *string `json:"deleted_at"`
}

type TransactionResponseMonthStatusSuccess struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalAmount  int    `json:"total_amount"`
	TotalSuccess int    `json:"total_success"`
}

type TransactionResponseYearStatusSuccess struct {
	Year         string `json:"year"`
	TotalAmount  int    `json:"total_amount"`
	TotalSuccess int    `json:"total_success"`
}

type TransactionResponseMonthStatusFailed struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
}

type TransactionResponseYearStatusFailed struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
	TotalFailed int    `json:"total_failed"`
}

type TransactionMonthMethodResponse struct {
	Month             string `json:"month"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int    `json:"total_transactions"`
	TotalAmount       int    `json:"total_amount"`
}

type TransactionYearMethodResponse struct {
	Year              string `json:"year"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int    `json:"total_transactions"`
	TotalAmount       int    `json:"total_amount"`
}

type TransactionMonthAmountResponse struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type TransactionYearlyAmountResponse struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}

type ApiResponseTransactionMonthStatusSuccess struct {
	Status  string                                   `json:"status"`
	Message string                                   `json:"message"`
	Data    []*TransactionResponseMonthStatusSuccess `json:"data"`
}

type ApiResponseTransactionYearStatusSuccess struct {
	Status  string                                  `json:"status"`
	Message string                                  `json:"message"`
	Data    []*TransactionResponseYearStatusSuccess `json:"data"`
}

type ApiResponseTransactionMonthStatusFailed struct {
	Status  string                                  `json:"status"`
	Message string                                  `json:"message"`
	Data    []*TransactionResponseMonthStatusFailed `json:"data"`
}

type ApiResponseTransactionYearStatusFailed struct {
	Status  string                                 `json:"status"`
	Message string                                 `json:"message"`
	Data    []*TransactionResponseYearStatusFailed `json:"data"`
}

type ApiResponseTransactionMonthMethod struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    []*TransactionMonthMethodResponse `json:"data"`
}

type ApiResponseTransactionYearMethod struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Data    []*TransactionYearMethodResponse `json:"data"`
}

type ApiResponseTransactionMonthAmount struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    []*TransactionMonthAmountResponse `json:"data"`
}

type ApiResponseTransactionYearAmount struct {
	Status  string                             `json:"status"`
	Message string                             `json:"message"`
	Data    []*TransactionYearlyAmountResponse `json:"data"`
}

type ApiResponseTransaction struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    *TransactionResponse `json:"data"`
}

type ApiResponseTransactions struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []*TransactionResponse `json:"data"`
}

type ApiResponseTransactionDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseTransactionAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationTransaction struct {
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Data       []*TransactionResponse `json:"data"`
	Pagination *PaginationMeta        `json:"pagination"`
}

type ApiResponsePaginationTransactionDeleteAt struct {
	Status     string                         `json:"status"`
	Message    string                         `json:"message"`
	Data       []*TransactionResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta                `json:"pagination"`
}
