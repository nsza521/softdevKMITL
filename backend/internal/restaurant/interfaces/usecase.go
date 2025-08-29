package interfaces

import (

)

type CustomerUsecase interface {
	Register(username string) error
	Login(username string, password string) (string, error)
}