package requests

import (
	methodtopup "MamangRust/paymentgatewaygrpc/pkg/method_topup"
	"errors"

	"github.com/go-playground/validator/v10"
)

type MonthTopupStatus struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

type MonthTopupStatusCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
	Month      int    `json:"month" validate:"required"`
}

type YearTopupStatusCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
}

type YearMonthMethod struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
}

type FindAllTopups struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllTopupsByCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Search     string `json:"search" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateTopupRequest struct {
	CardNumber  string `json:"card_number" validate:"required,min=1"`
	TopupAmount int    `json:"topup_amount" validate:"required,min=50000"`
	TopupMethod string `json:"topup_method" validate:"required"`
}

type UpdateTopupRequest struct {
	CardNumber  string `json:"card_number" validate:"required,min=1"`
	TopupID     *int   `json:"topup_id"`
	TopupAmount int    `json:"topup_amount" validate:"required,min=50000"`
	TopupMethod string `json:"topup_method" validate:"required"`
}

type UpdateTopupAmount struct {
	TopupID     int `json:"topup_id" validate:"required,min=1"`
	TopupAmount int `json:"topup_amount" validate:"required,min=50000"`
}

type UpdateTopupStatus struct {
	TopupID int    `json:"topup_id" validate:"required,min=1"`
	Status  string `json:"status" validate:"required"`
}

func (r *CreateTopupRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.TopupAmount < 50000 {
		return errors.New("topup amount must be greater than or equal to 50000")
	}

	if r.TopupMethod == "" {
		return errors.New("top-up method is required")
	}

	if !methodtopup.PaymentMethodValidator(r.TopupMethod) {
		return errors.New("topup method not found")
	}

	return nil
}

func (r *UpdateTopupRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	if *r.TopupID <= 0 {
		return errors.New("top-up ID must be a positive integer")
	}

	if r.TopupAmount < 50000 {
		return errors.New("topup amount must be greater than or equal to 50000")
	}

	if r.TopupMethod == "" {
		return errors.New("top-up method is required")
	}

	if !methodtopup.PaymentMethodValidator(r.TopupMethod) {
		return errors.New("topup method not found")
	}

	return nil
}

func (r *UpdateTopupAmount) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.TopupID <= 0 {
		return errors.New("top-up ID must be a positive integer")
	}

	if r.TopupAmount < 50000 {
		return errors.New("topup amount must be greater than or equal to 50000")
	}

	return nil
}

func (r *UpdateTopupStatus) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
