package record

type CardRecord struct {
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

type CardMonthBalance struct {
	Month        string `json:"month"`
	TotalBalance int64  `json:"total_balance"`
}

type CardYearlyBalance struct {
	Year         string `json:"year"`
	TotalBalance int64  `json:"total_balance"`
}

type CardMonthAmount struct {
	Month       string `json:"month"`
	TotalAmount int64  `json:"total_amount"`
}

type CardYearAmount struct {
	Year        string `json:"year"`
	TotalAmount int64  `json:"total_amount"`
}


