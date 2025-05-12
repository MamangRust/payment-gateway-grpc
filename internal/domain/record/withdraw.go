package record

type WithdrawRecord struct {
	ID             int     `json:"id"`
	WithdrawNo     string  `json:"withdraw_no"`
	CardNumber     string  `json:"card_number"`
	WithdrawAmount int     `json:"withdraw_amount"`
	WithdrawTime   string  `json:"withdraw_time"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}

type WithdrawRecordMonthStatusSuccess struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type WithdrawRecordYearStatusSuccess struct {
	Year         string `json:"year"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type WithdrawRecordMonthStatusFailed struct {
	Year        string `json:"year"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type WithdrawRecordYearStatusFailed struct {
	Year        string `json:"year"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type WithdrawMonthlyAmount struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type WithdrawYearlyAmount struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}
