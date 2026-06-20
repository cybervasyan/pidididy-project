package model

import "errors"

var (
	ErrorOrderNotInPending = errors.New("Заказ не в статусе PENDING_PAYMENT")
	ErrorPartNotFound      = errors.New("Такой запчасти не существует")
)
