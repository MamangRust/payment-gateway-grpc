package requests

import (
	methodtopup "MamangRust/paymentgatewaygrpc/pkg/method_topup"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type MonthYearCardNumberCard struct {
	CardNumber string `json:"card_number" validate:"required"`
	Year       int    `json:"year" validate:"required"`
}

type FindAllCards struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateCardRequest struct {
	UserID       int       `json:"user_id"`
	CardType     string    `json:"card_type" validate:"required"`
	ExpireDate   time.Time `json:"expire_date" validate:"required"`
	CVV          string    `json:"cvv" validate:"required"`
	CardProvider string    `json:"card_provider" validate:"required"`
}

func (r *CreateCardRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if r.CardType != "credit" && r.CardType != "debit" {
		return fmt.Errorf("card type must be credit or debit")
	}

	if !methodtopup.PaymentMethodValidator(r.CardProvider) {
		return fmt.Errorf("card provider not found")
	}

	if err != nil {
		return err
	}

	return nil
}

type UpdateCardRequest struct {
	CardID       int       `json:"card_id" validate:"required,min=1"`
	UserID       int       `json:"user_id" validate:"required,min=1"`
	CardType     string    `json:"card_type" validate:"required"`
	ExpireDate   time.Time `json:"expire_date" validate:"required"`
	CVV          string    `json:"cvv" validate:"required"`
	CardProvider string    `json:"card_provider" validate:"required"`
}

func (r *UpdateCardRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if r.CardType != "credit" && r.CardType != "debit" {
		return fmt.Errorf("card type must be credit or debit")
	}

	if !methodtopup.PaymentMethodValidator(r.CardProvider) {
		return fmt.Errorf("card provider not found")
	}

	if err != nil {
		return err
	}

	return nil
}
