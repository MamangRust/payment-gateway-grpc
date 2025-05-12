package record

type TransferRecord struct {
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

type TransferRecordMonthStatusSuccess struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type TransferRecordYearStatusSuccess struct {
	Year         string `json:"year"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type TransferRecordMonthStatusFailed struct {
	Year        string `json:"year"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type TransferRecordYearStatusFailed struct {
	Year        string `json:"year"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type TransferMonthAmount struct {
	Month       string `json:"month"`
	TotalAmount int    `json:"total_amount"`
}

type TransferYearAmount struct {
	Year        string `json:"year"`
	TotalAmount int    `json:"total_amount"`
}
