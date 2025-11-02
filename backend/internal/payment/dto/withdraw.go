package dto

import (
	// "github.com/google/uuid"
)

type WithdrawRequest struct {
	FullName          string  `json:"full_name" binding:"required"`
	BankName          string  `json:"bank_name" binding:"required"`
	BankAccountNumber string  `json:"bank_account_number" binding:"required,numeric,min=10,max=12"`
	WithdrawAmount            float32 `json:"withdraw_amount" binding:"required,gt=0"`
}

type WithdrawResponse struct {
	RemainingBalance float32 `json:"remaining_balance"`
}