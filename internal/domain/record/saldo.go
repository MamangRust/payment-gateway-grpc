package record

type SaldoRecord struct {
	ID             int     `json:"id"`
	CardNumber     string  `json:"card_number"`
	TotalBalance   int     `json:"total_balance"`
	WithdrawAmount int     `json:"withdraw_amount"`
	WithdrawTime   string  `json:"withdraw_time"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}

type SaldoMonthTotalBalance struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalBalance int    `json:"total_balance"`
}

type SaldoYearTotalBalance struct {
	Year         string `json:"year"`
	TotalBalance int    `json:"total_balance"`
}

type SaldoMonthSaldoBalance struct {
	Month        string `json:"month"`
	TotalBalance int    `json:"total_balance"`
}

type SaldoYearSaldoBalance struct {
	Year         string `json:"year"`
	TotalBalance int    `json:"total_balance"`
}
