package model

import "errors"

var (
	ErrorOrderDoesntExist  = errors.New("Такого заказа не существует")
	ErrorOrderNotInPending = errors.New("Заказ не в статусе PENDING_PAYMENT")
)
