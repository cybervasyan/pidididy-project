package converter

import (
	"github.com/cybervasyan/pdididy-project/order/internal/model"
	orderv1 "github.com/cybervasyan/pdididy-project/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func ModelToOrderDto(order model.Order) *orderv1.OrderDto {
	dto := &orderv1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     orderv1.OrderStatus(order.Status),
	}

	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderv1.NewOptPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	return dto
}

func CreateRequestToModel(req *orderv1.CreateOrderRequest) model.Order {
	return model.Order{
		UserUUID:  req.UserUUID,
		PartUuids: req.PartUuids,
	}
}

func ModelToCreateResponse(order model.Order) *orderv1.CreateOrderResponse {
	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}
}

func PaymentMethodToModel(method orderv1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func PayRequestToModel(orderUUID uuid.UUID, req *orderv1.PayOrderRequest) model.Order {
	pm := PaymentMethodToModel(req.PaymentMethod)
	return model.Order{
		OrderUUID:     orderUUID,
		PaymentMethod: &pm,
	}
}

func ModelToPayResponse(order model.Order) *orderv1.PayOrderResponse {
	var transactionUUID uuid.UUID
	if order.TransactionUUID != nil {
		transactionUUID = *order.TransactionUUID
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}
}
