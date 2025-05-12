package response

type TopupResponse struct {
	ID          int    `json:"id"`
	CardNumber  string `json:"card_number"`
	TopupNo     string `json:"topup_no"`
	TopupAmount int    `json:"topup_amount"`
	TopupMethod string `json:"topup_method"`
	TopupTime   string `json:"topup_time"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TopupResponseDeleteAt struct {
	ID          int     `json:"id"`
	CardNumber  string  `json:"card_number"`
	TopupNo     string  `json:"topup_no"`
	TopupAmount int     `json:"topup_amount"`
	TopupMethod string  `json:"topup_method"`
	TopupTime   string  `json:"topup_time"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_At"`
}

type TopupResponseMonthStatusSuccess struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalAmount  int    `json:"total_amount"`
	TotalSuccess int    `json:"total_success"`
}

type TopupResponseYearStatusSuccess struct {
	Year         string `json:"year"`
	TotalAmount  int    `json:"total_amount"`
	TotalSuccess int    `json:"total_success"`
}

type TopupResponseMonthStatusFailed struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
}

type TopupResponseYearStatusFailed struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
	TotalFailed int    `json:"total_failed"`
}

type TopupMonthMethodResponse struct {
	Month       string `json:"month"`
	TopupMethod string `json:"topup_method"`
	TotalTopups int    `json:"total_topups"`
	TotalAmount int    `json:"total_amount"`
}

type TopupYearlyMethodResponse struct {
	Year        string `json:"year"`
	TopupMethod string `json:"topup_method"`
	TotalTopups int    `json:"total_topups"`
	TotalAmount int    `json:"total_amount"`
}

type TopupMonthAmountResponse struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type TopupYearlyAmountResponse struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}

type ApiResponseTopupMonthStatusSuccess struct {
	Status  string                             `json:"status"`
	Message string                             `json:"message"`
	Data    []*TopupResponseMonthStatusSuccess `json:"data"`
}

type ApiResponseTopupYearStatusSuccess struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    []*TopupResponseYearStatusSuccess `json:"data"`
}

type ApiResponseTopupMonthStatusFailed struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    []*TopupResponseMonthStatusFailed `json:"data"`
}

type ApiResponseTopupYearStatusFailed struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Data    []*TopupResponseYearStatusFailed `json:"data"`
}

type ApiResponseTopupMonthMethod struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []*TopupMonthMethodResponse `json:"data"`
}

type ApiResponseTopupYearMethod struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []*TopupYearlyMethodResponse `json:"data"`
}

type ApiResponseTopupMonthAmount struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []*TopupMonthAmountResponse `json:"data"`
}

type ApiResponseTopupYearAmount struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []*TopupYearlyAmountResponse `json:"data"`
}

type ApiResponseTopup struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    *TopupResponse `json:"data"`
}

type ApiResponseTopupDeleteAt struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    *TopupResponseDeleteAt `json:"data"`
}

type ApiResponsesTopup struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    []*TopupResponse `json:"data"`
}

type ApiResponsePaginationTopup struct {
	Status     string           `json:"status"`
	Message    string           `json:"message"`
	Data       []*TopupResponse `json:"data"`
	Pagination *PaginationMeta  `json:"pagination"`
}

type ApiResponsePaginationTopupDeleteAt struct {
	Status     string                   `json:"status"`
	Message    string                   `json:"message"`
	Data       []*TopupResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta          `json:"pagination"`
}

type ApiResponseTopupDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseTopupAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
