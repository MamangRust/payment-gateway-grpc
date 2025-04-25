package response

type SaldoResponse struct {
	ID             int    `json:"id"`
	CardNumber     string `json:"card_number"`
	TotalBalance   int    `json:"total_balance"`
	WithdrawAmount int    `json:"withdraw_amount"`
	WithdrawTime   string `json:"withdraw_time"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type SaldoResponseDeleteAt struct {
	ID             int     `json:"id"`
	CardNumber     string  `json:"card_number"`
	TotalBalance   int     `json:"total_balance"`
	WithdrawAmount int     `json:"withdraw_amount"`
	WithdrawTime   string  `json:"withdraw_time"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}

type SaldoMonthTotalBalanceResponse struct {
	Month        string `json:"month"`
	Year         string `json:"year"`
	TotalBalance int    `json:"total_balance"`
}

type SaldoYearTotalBalanceResponse struct {
	Year         string `json:"year"`
	TotalBalance int    `json:"total_balance"`
}

type SaldoMonthBalanceResponse struct {
	Month        string `json:"month"`
	TotalBalance int    `json:"total_balance"`
}

type SaldoYearBalanceResponse struct {
	Year         string `json:"year"`
	TotalBalance int    `json:"total_balance"`
}

type ApiResponseSaldo struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    SaldoResponse `json:"data"`
}

type ApiResponsesSaldo struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    []*SaldoResponse `json:"data"`
}

type ApiResponseSaldoDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseSaldoAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseMonthTotalSaldo struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    []*SaldoMonthTotalBalanceResponse `json:"data"`
}

type ApiResponseYearTotalSaldo struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Data    []*SaldoYearTotalBalanceResponse `json:"data"`
}

type ApiResponseMonthSaldoBalances struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []*SaldoMonthBalanceResponse `json:"data"`
}

type ApiResponseYearSaldoBalances struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []*SaldoYearBalanceResponse `json:"data"`
}

type ApiResponsePaginationSaldo struct {
	Status     string           `json:"status"`
	Message    string           `json:"message"`
	Data       []*SaldoResponse `json:"data"`
	Pagination *PaginationMeta  `json:"pagination"`
}

type ApiResponsePaginationSaldoDeleteAt struct {
	Status     string                   `json:"status"`
	Message    string                   `json:"message"`
	Data       []*SaldoResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta          `json:"pagination"`
}
