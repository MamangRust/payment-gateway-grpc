package requests

import "github.com/go-playground/validator/v10"

type MonthYearPaymentMethodApiKey struct {
	Apikey string `json:"api_key" validate:"required,min=1"`
	Year   int    `json:"year" validate:"required"`
}

type MonthYearAmountApiKey struct {
	Apikey string `json:"api_key" validate:"required,min=1"`
	Year   int    `json:"year" validate:"required"`
}

type MonthYearTotalAmountApiKey struct {
	Apikey string `json:"api_key" validate:"required,min=1"`
	Year   int    `json:"year" validate:"required"`
}

type MonthYearPaymentMethodMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required,min=1"`
	Year       int `json:"year" validate:"required"`
}

type MonthYearAmountMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required,min=1"`
	Year       int `json:"year" validate:"required"`
}

type MonthYearTotalAmountMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required,min=1"`
	Year       int `json:"year" validate:"required"`
}

type FindAllMerchants struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllMerchantTransactions struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllMerchantTransactionsById struct {
	MerchantID int    `json:"merchant_id" validate:"required,min=1"`
	Search     string `json:"search" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllMerchantTransactionsByApiKey struct {
	ApiKey   string `json:"api_key" validate:"required"`
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateMerchantRequest struct {
	Name   string `json:"name" validate:"required"`
	UserID int    `json:"user_id" validate:"required,min=1"`
}

type UpdateMerchantRequest struct {
	MerchantID *int   `json:"merchant_id"`
	Name       string `json:"name" validate:"required"`
	UserID     int    `json:"user_id" validate:"required,min=1"`
	Status     string `json:"status" validate:"required"`
}

type UpdateMerchantStatus struct {
	MerchantID int    `json:"merchant_id" validate:"required,min=1"`
	Status     string `json:"status" validate:"required"`
}

func (r CreateMerchantRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

func (r UpdateMerchantRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

func (r UpdateMerchantStatus) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}
