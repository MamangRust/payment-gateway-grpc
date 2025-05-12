package record

type TopupRecord struct {
	ID          int     `json:"id"`
	CardNumber  string  `json:"card_number"`
	TopupNo     string  `json:"topup_no"`
	TopupAmount int     `json:"topup_amount"`
	TopupMethod string  `json:"topup_method"`
	TopupTime   string  `json:"topup_time"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at"`
}

type TopupRecordMonthStatusSuccess struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type TopupRecordYearStatusSuccess struct {
	Year         string `json:"year"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type TopupRecordMonthStatusFailed struct {
	Year        string `json:"year"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type TopupRecordYearStatusFailed struct {
	Year        string `json:"year"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type TopupMonthMethod struct {
	Month       string `json:"month"`
	TopupMethod string `json:"topup_method"`
	TotalTopups int    `json:"total_topups"`
	TotalAmount int    `json:"total_amount"`
}

type TopupYearlyMethod struct {
	Year        string `json:"year"`
	TopupMethod string `json:"topup_method"`
	TotalTopups int    `json:"total_topups"`
	TotalAmount int    `json:"total_amount"`
}

type TopupMonthAmount struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type TopupYearlyAmount struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}
