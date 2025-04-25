package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type MonthYearCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
}

type MonthStatusTransferCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
	Month      int    `json:"month" validate:"required"`
}

type YearStatusTransferCardNumber struct {
	CardNumber string `json:"card_number" validate:"required,min=1"`
	Year       int    `json:"year" validate:"required"`
}

type MonthStatusTransfer struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

type FindAllTranfers struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateTransferRequest struct {
	TransferFrom   string `json:"transfer_from" validate:"required"`
	TransferTo     string `json:"transfer_to" validate:"required,min=1"`
	TransferAmount int    `json:"transfer_amount" validate:"required,min=50000"`
}

type UpdateTransferRequest struct {
	TransferID     *int   `json:"transfer_id"`
	TransferFrom   string `json:"transfer_from" validate:"required"`
	TransferTo     string `json:"transfer_to" validate:"required,min=1"`
	TransferAmount int    `json:"transfer_amount" validate:"required,min=50000"`
}

type UpdateTransferAmountRequest struct {
	TransferID     int `json:"transfer_id" validate:"required,min=1"`
	TransferAmount int `json:"transfer_amount" validate:"required,gt=0"`
}

type UpdateTransferStatus struct {
	TransferID int    `json:"transfer_id" validate:"required,min=1"`
	Status     string `json:"status" validate:"required"`
}

func (r *CreateTransferRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.TransferAmount < 50000 {
		return errors.New("transfer amount must be at least 50,000")
	}

	return nil
}

func (r *UpdateTransferRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if *r.TransferID <= 0 {
		return errors.New("transfer ID must be a positive integer")
	}

	if r.TransferAmount < 50000 {
		return errors.New("transfer amount must be at least 50,000")
	}

	return nil
}

func (r *UpdateTransferAmountRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.TransferID <= 0 {
		return errors.New("transfer ID must be a positive integer")
	}

	if r.TransferAmount <= 0 {
		return errors.New("transfer amount must be greater than zero")
	}

	return nil
}

func (r *UpdateTransferStatus) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
