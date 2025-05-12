package response

type TransferResponse struct {
	ID             int    `json:"id"`
	TransferNo     string `json:"transfer_no"`
	TransferFrom   string `json:"transfer_from"`
	TransferTo     string `json:"transfer_to"`
	TransferAmount int    `json:"transfer_amount"`
	TransferTime   string `json:"transfer_time"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type TransferResponseDeleteAt struct {
	ID             int     `json:"id"`
	TransferNo     string  `json:"transfer_no"`
	TransferFrom   string  `json:"transfer_from"`
	TransferTo     string  `json:"transfer_to"`
	TransferAmount int     `json:"transfer_amount"`
	TransferTime   string  `json:"transfer_time"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}

type TransferResponseMonthStatusSuccess struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalAmount  int    `json:"total_amount"`
	TotalSuccess int    `json:"total_success"`
}

type TransferResponseYearStatusSuccess struct {
	Year         string `json:"year"`
	TotalAmount  int    `json:"total_amount"`
	TotalSuccess int    `json:"total_success"`
}

type TransferResponseMonthStatusFailed struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
}

type TransferResponseYearStatusFailed struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
	TotalFailed int    `json:"total_failed"`
}

type TransferMonthAmountResponse struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type TransferYearAmountResponse struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}

type ApiResponseTransferMonthStatusSuccess struct {
	Status  string                                `json:"status"`
	Message string                                `json:"message"`
	Data    []*TransferResponseMonthStatusSuccess `json:"data"`
}

type ApiResponseTransferYearStatusSuccess struct {
	Status  string                               `json:"status"`
	Message string                               `json:"message"`
	Data    []*TransferResponseYearStatusSuccess `json:"data"`
}

type ApiResponseTransferMonthStatusFailed struct {
	Status  string                               `json:"status"`
	Message string                               `json:"message"`
	Data    []*TransferResponseMonthStatusFailed `json:"data"`
}

type ApiResponseTransferYearStatusFailed struct {
	Status  string                              `json:"status"`
	Message string                              `json:"message"`
	Data    []*TransferResponseYearStatusFailed `json:"data"`
}

type ApiResponseTransferMonthAmount struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Data    []*TransferMonthAmountResponse `json:"data"`
}

type ApiResponseTransferYearAmount struct {
	Status  string                        `json:"status"`
	Message string                        `json:"message"`
	Data    []*TransferYearAmountResponse `json:"data"`
}

type ApiResponseTransfer struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *TransferResponse `json:"data"`
}

type ApiResponseTransfers struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []*TransferResponse `json:"data"`
}

type ApiResponseTransferDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseTransferAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationTransfer struct {
	Status     string              `json:"status"`
	Message    string              `json:"message"`
	Data       []*TransferResponse `json:"data"`
	Pagination *PaginationMeta     `json:"pagination"`
}

type ApiResponsePaginationTransferDeleteAt struct {
	Status     string                      `json:"status"`
	Message    string                      `json:"message"`
	Data       []*TransferResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta             `json:"pagination"`
}
