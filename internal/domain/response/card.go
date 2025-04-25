package response

type CardResponse struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	CardNumber   string `json:"card_number"`
	CardType     string `json:"card_type"`
	ExpireDate   string `json:"expire_date"`
	CVV          string `json:"cvv"`
	CardProvider string `json:"card_provider"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CardResponseDeleteAt struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	CardNumber   string  `json:"card_number"`
	CardType     string  `json:"card_type"`
	ExpireDate   string  `json:"expire_date"`
	CVV          string  `json:"cvv"`
	CardProvider string  `json:"card_provider"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    *string `json:"deleted_at"`
}

type DashboardCard struct {
	TotalBalance     *int64 `json:"total_balance"`
	TotalTopup       *int64 `json:"total_topup"`
	TotalWithdraw    *int64 `json:"total_withdraw"`
	TotalTransaction *int64 `json:"total_transaction"`
	TotalTransfer    *int64 `json:"total_transfer"`
}

type DashboardCardCardNumber struct {
	TotalBalance          *int64 `json:"total_balance"`
	TotalTopup            *int64 `json:"total_topup"`
	TotalWithdraw         *int64 `json:"total_withdraw"`
	TotalTransaction      *int64 `json:"total_transaction"`
	TotalTransferSend     *int64 `json:"total_transfer_send"`
	TotalTransferReceiver *int64 `json:"total_transfer_receiver"`
}

type CardResponseMonthBalance struct {
	Month        string `json:"month"`
	TotalBalance int64  `json:"total_balance"`
}

type CardResponseYearlyBalance struct {
	Year         string `json:"year"`
	TotalBalance int64  `json:"total_balance"`
}

type CardResponseMonthAmount struct {
	Month       string `json:"month"`
	TotalAmount int64  `json:"total_amount"`
}

type CardResponseYearAmount struct {
	Year        string `json:"year"`
	TotalAmount int64  `json:"total_amount"`
}

type ApiResponseCard struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *CardResponse `json:"data"`
}

type ApiResponseCardDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseCardAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationCard struct {
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Data       []*CardResponse `json:"data"`
	Pagination *PaginationMeta `json:"pagination"`
}

type ApiResponsePaginationCardDeleteAt struct {
	Status     string                  `json:"status"`
	Message    string                  `json:"message"`
	Data       []*CardResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta         `json:"pagination"`
}

type ApiResponseMonthlyBalance struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []*CardResponseMonthBalance `json:"data"`
}

type ApiResponseYearlyBalance struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []*CardResponseYearlyBalance `json:"data"`
}

type ApiResponseMonthlyAmount struct {
	Status  string                     `json:"status"`
	Message string                     `json:"message"`
	Data    []*CardResponseMonthAmount `json:"data"`
}

type ApiResponseYearlyAmount struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    []*CardResponseYearAmount `json:"data"`
}

type ApiResponseDashboardCard struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    *DashboardCard `json:"data"`
}

type ApiResponseDashboardCardNumber struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    *DashboardCardCardNumber `json:"data"`
}
