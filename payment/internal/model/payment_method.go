package model

// PaymentMethod — способ оплаты в терминах домена. Типизированная строка,
// а не голый string: компилятор не даст подсунуть произвольное значение,
// а набор валидных вариантов задан константами ниже.
type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "PAYMENT_METHOD_UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "PAYMENT_METHOD_CARD"
	PaymentMethodSBP           PaymentMethod = "PAYMENT_METHOD_SBP"
	PaymentMethodCreditCard    PaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)
