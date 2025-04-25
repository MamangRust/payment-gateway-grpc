package requests

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type YearMonthCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
}

type MonthStatusWithdraw struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

type MonthStatusWithdrawCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
	Month      int    `json:"month" validate:"required"`
}

type YearStatusWithdrawCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
}

type FindAllWithdraws struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllWithdrawCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Search     string `json:"search" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateWithdrawRequest struct {
	CardNumber     string    `json:"card_number" validate:"required,min=1"`
	WithdrawAmount int       `json:"withdraw_amount" validate:"required,min=50000"`
	WithdrawTime   time.Time `json:"withdraw_time" validate:"required"`
}

type UpdateWithdrawRequest struct {
	CardNumber     string    `json:"card_number" validate:"required,min=1"`
	WithdrawID     *int      `json:"withdraw_id"`
	WithdrawAmount int       `json:"withdraw_amount" validate:"required,min=50000"`
	WithdrawTime   time.Time `json:"withdraw_time" validate:"required"`
}

type UpdateWithdrawStatus struct {
	WithdrawID int    `json:"withdraw_id" validate:"required,min=1"`
	Status     string `json:"status" validate:"required"`
}

func (r *CreateWithdrawRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.WithdrawAmount < 50000 {
		return errors.New("withdraw amount must be at least 50,000")
	}

	if r.WithdrawTime.After(time.Now()) {
		return errors.New("withdraw time cannot be in the future")
	}

	return nil
}

func (r *UpdateWithdrawRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.WithdrawAmount < 50000 {
		return errors.New("withdraw amount must be at least 50,000")
	}

	if r.WithdrawTime.After(time.Now()) {
		return errors.New("withdraw time cannot be in the future")
	}

	return nil
}

func (r *UpdateWithdrawStatus) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
