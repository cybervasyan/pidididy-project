package model

import "errors"

var (
	ErrOrderUUIDRequired        = errors.New("order uuid обязателен")
	ErrUserUUIDRequired         = errors.New("user uuid обязателен")
	ErrPaymentMethodUnspecified = errors.New("способ оплаты не указан")
)
