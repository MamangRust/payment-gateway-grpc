package requests

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type FindAllSaldos struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type MonthTotalSaldoBalance struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

type CreateSaldoRequest struct {
	CardNumber   string `json:"card_number" validate:"required"`
	TotalBalance int    `json:"total_balance" validate:"required"`
}

type UpdateSaldoRequest struct {
	SaldoID      *int   `json:"saldo_id" `
	CardNumber   string `json:"card_number" validate:"required"`
	TotalBalance int    `json:"total_balance" validate:"required"`
}

type UpdateSaldoBalance struct {
	CardNumber   string `json:"card_number" validate:"required,min=1"`
	TotalBalance int    `json:"total_balance" validate:"required,min=50000"`
}

type UpdateSaldoWithdraw struct {
	CardNumber     string     `json:"card_number" validate:"required,min=1"`
	TotalBalance   int        `json:"total_balance" validate:"required,min=50000"`
	WithdrawAmount *int       `json:"withdraw_amount" validate:"omitempty,gte=0"`
	WithdrawTime   *time.Time `json:"withdraw_time" validate:"omitempty"`
}

func (r *CreateSaldoRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateSaldoRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateSaldoBalance) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateSaldoWithdraw) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.WithdrawAmount != nil && r.WithdrawTime == nil {
		return errors.New("withdraw time must be provided if withdraw amount is provided")
	}
	if r.WithdrawAmount == nil && r.WithdrawTime != nil {
		return errors.New("withdraw amount must be provided if withdraw time is provided")
	}
	if r.WithdrawAmount != nil && *r.WithdrawAmount > r.TotalBalance {
		return errors.New("withdraw amount cannot be greater than total balance")
	}
	return nil
}
