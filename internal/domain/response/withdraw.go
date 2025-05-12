package response

type WithdrawResponse struct {
	ID             int    `json:"id"`
	WithdrawNo     string `json:"withdraw_no"`
	CardNumber     string `json:"card_number"`
	WithdrawAmount int    `json:"withdraw_amount"`
	WithdrawTime   string `json:"withdraw_time"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type WithdrawResponseDeleteAt struct {
	ID             int     `json:"id"`
	WithdrawNo     string  `json:"withdraw_no"`
	CardNumber     string  `json:"card_number"`
	WithdrawAmount int     `json:"withdraw_amount"`
	WithdrawTime   string  `json:"withdraw_time"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}

type WithdrawResponseMonthStatusSuccess struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalAmount  int    `json:"total_amount"`
	TotalSuccess int    `json:"total_success"`
}

type WithdrawResponseYearStatusSuccess struct {
	Year         string `json:"year"`
	TotalAmount  int    `json:"total_amount"`
	TotalSuccess int    `json:"total_success"`
}

type WithdrawResponseMonthStatusFailed struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
}

type WithdrawResponseYearStatusFailed struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
	TotalFailed int    `json:"total_failed"`
}

type WithdrawMonthlyAmountResponse struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type WithdrawYearlyAmountResponse struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}

type ApiResponseWithdrawMonthStatusSuccess struct {
	Status  string                                `json:"status"`
	Message string                                `json:"message"`
	Data    []*WithdrawResponseMonthStatusSuccess `json:"data"`
}

type ApiResponseWithdrawYearStatusSuccess struct {
	Status  string                               `json:"status"`
	Message string                               `json:"message"`
	Data    []*WithdrawResponseYearStatusSuccess `json:"data"`
}

type ApiResponseWithdrawMonthStatusFailed struct {
	Status  string                               `json:"status"`
	Message string                               `json:"message"`
	Data    []*WithdrawResponseMonthStatusFailed `json:"data"`
}

type ApiResponseWithdrawYearStatusFailed struct {
	Status  string                              `json:"status"`
	Message string                              `json:"message"`
	Data    []*WithdrawResponseYearStatusFailed `json:"data"`
}

type ApiResponseWithdrawMonthAmount struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Data    []*WithdrawMonthlyAmountResponse `json:"data"`
}

type ApiResponseWithdrawYearAmount struct {
	Status  string                          `json:"status"`
	Message string                          `json:"message"`
	Data    []*WithdrawYearlyAmountResponse `json:"data"`
}

type ApiResponsesWithdraw struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []*WithdrawResponse `json:"data"`
}

type ApiResponseWithdraw struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *WithdrawResponse `json:"data"`
}

type ApiResponseWithdrawDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseWithdrawAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationWithdraw struct {
	Status     string              `json:"status"`
	Message    string              `json:"message"`
	Data       []*WithdrawResponse `json:"data"`
	Pagination PaginationMeta      `json:"pagination"`
}

type ApiResponsePaginationWithdrawDeleteAt struct {
	Status     string                      `json:"status"`
	Message    string                      `json:"message"`
	Data       []*WithdrawResponseDeleteAt `json:"data"`
	Pagination PaginationMeta              `json:"pagination"`
}
