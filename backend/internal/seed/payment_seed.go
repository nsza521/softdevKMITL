package seed

import (
	// "fmt"
	// "time"

	"gorm.io/gorm"
	
	"backend/internal/db_model"
)

func seedPaymentMethods(db *gorm.DB) error {
	var paymentMethods []models.PaymentMethod
	paymentMethods = []models.PaymentMethod{
		{
			Name:        "Promtpay",
			Type:        "topup",
			Description: "PromptPay is a Thai electronic payment system that allows users to make instant money transfers using a mobile number or national ID.",
			ImageURL:    nil,
		},
		{
			Name:        "KBANK",
			Type:        "topup",
			Description: "Kasikornbank (KBank) is one of the largest banks in Thailand, offering a wide range of financial services including personal and business banking, loans, and investment products.",
			ImageURL:    nil,
		},
		{
			Name:        "SCB",
			Type:        "topup",
			Description: "Siam Commercial Bank (SCB) is a leading bank in Thailand, providing various banking services such as savings accounts, loans, credit cards, and investment options.",
			ImageURL:    nil,
		},
		{
			Name:        "Wallet",
			Type:        "paid",
			Description: "A digital wallet is a software-based system that securely stores users' payment information and passwords for numerous payment methods and websites.",
			ImageURL:    nil,
		},
	}

	for _, method := range paymentMethods {
		if err := db.Where("name = ?", method.Name).FirstOrCreate(&method).Error; err != nil {
			return err
		}
	}

	return nil
}