package dto

type RegisterRestaurantRequest struct {
	Username      	string `json:"username" binding:"required"`
	Email       	string `json:"email" binding:"required,email"`
	Password    	string `json:"password" binding:"required,min=8"`
	BankName    	string `json:"bank_name" binding:"required"`
	AccountNumber 	string `json:"account_no" binding:"required"`
	AccountName   	string `json:"account_name" binding:"required"`
	Name		 	string `json:"name" binding:"required"`
}