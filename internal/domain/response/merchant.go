package response

type MerchantResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserID    int    `json:"user_id"`
	ApiKey    string `json:"api_key"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MerchantResponseDeleteAt struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	UserID    int     `json:"user_id"`
	ApiKey    string  `json:"api_key"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type MerchantTransactionResponse struct {
	ID              int    `json:"id"`
	CardNumber      string `json:"card_number"`
	Amount          int32  `json:"amount"`
	PaymentMethod   string `json:"payment_method"`
	MerchantID      int32  `json:"merchant_id"`
	MerchantName    string `json:"merchant_name"`
	TransactionTime string `json:"transaction_time"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type MerchantResponseMonthlyPaymentMethod struct {
	Month         string `json:"month"`
	PaymentMethod string `json:"payment_method"`
	TotalAmount   int    `json:"total_amount"`
}

type MerchantResponseYearlyPaymentMethod struct {
	Year          string `json:"year"`
	PaymentMethod string `json:"payment_method"`
	TotalAmount   int    `json:"total_amount"`
}

type MerchantResponseMonthlyAmount struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type MerchantResponseYearlyAmount struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}

type MerchantResponseMonthlyTotalAmount struct {
	Year        string `json:"year"`
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type MerchantResponseYearlyTotalAmount struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}

type ApiResponseMerchantMonthlyPaymentMethod struct {
	Status  string                                  `json:"status"`
	Message string                                  `json:"message"`
	Data    []*MerchantResponseMonthlyPaymentMethod `json:"data"`
}

type ApiResponseMerchantYearlyPaymentMethod struct {
	Status  string                                 `json:"status"`
	Message string                                 `json:"message"`
	Data    []*MerchantResponseYearlyPaymentMethod `json:"data"`
}

type ApiResponseMerchantMonthlyAmount struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Data    []*MerchantResponseMonthlyAmount `json:"data"`
}

type ApiResponseMerchantYearlyAmount struct {
	Status  string                          `json:"status"`
	Message string                          `json:"message"`
	Data    []*MerchantResponseYearlyAmount `json:"data"`
}

type ApiResponseMerchantMonthlyTotalAmount struct {
	Status  string                                `json:"status"`
	Message string                                `json:"message"`
	Data    []*MerchantResponseMonthlyTotalAmount `json:"data"`
}

type ApiResponseMerchantYearlyTotalAmount struct {
	Status  string                               `json:"status"`
	Message string                               `json:"message"`
	Data    []*MerchantResponseYearlyTotalAmount `json:"data"`
}

type ApiResponsesMerchant struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []*MerchantResponse `json:"data"`
}

type ApiResponseMerchant struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    MerchantResponse `json:"data"`
}

type ApiResponseMerchantDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseMerchantAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationMerchant struct {
	Status     string              `json:"status"`
	Message    string              `json:"message"`
	Data       []*MerchantResponse `json:"data"`
	Pagination *PaginationMeta     `json:"pagination"`
}

type ApiResponsePaginationMerchantDeleteAt struct {
	Status     string                      `json:"status"`
	Message    string                      `json:"message"`
	Data       []*MerchantResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta             `json:"pagination"`
}

type ApiResponsePaginationMerchantTransaction struct {
	Status     string                         `json:"status"`
	Message    string                         `json:"message"`
	Data       []*MerchantTransactionResponse `json:"data"`
	Pagination *PaginationMeta                `json:"pagination"`
}
