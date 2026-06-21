package model

import "errors"

var (
	ErrOrderNotInPending = errors.New("Заказ не в статусе PENDING_PAYMENT")
	ErrPartNotFound      = errors.New("Такой запчасти не существует")
	ErrOrderDoesntExist  = errors.New("Такого заказа не существует")
)
