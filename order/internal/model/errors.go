package model

import "errors"

var (
	ErrOrderNotInPending = errors.New("заказ не в статусе PENDING_PAYMENT")
	ErrPartNotFound      = errors.New("такой запчасти не существует")
	ErrOrderDoesntExist  = errors.New("такого заказа не существует")
)
