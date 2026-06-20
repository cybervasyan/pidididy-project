package converter

import (
	"github.com/cybervasyan/pdididy-project/order/internal/model"
	repoModel "github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func ToServiceModel(order repoModel.Order) model.Order {
	result := model.Order{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     model.OrderStatus(order.Status),
	}

	if order.PaymentMethod != nil {
		pm := model.PaymentMethod(*order.PaymentMethod)
		result.PaymentMethod = &pm
	}

	if order.TransactionUUID != nil {
		uid := *order.TransactionUUID
		result.TransactionUUID = &uid
	}

	return result
}

func ToRepoModel(order model.Order) repoModel.Order {
	result := repoModel.Order{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     repoModel.OrderStatus(order.Status),
	}

	if order.PaymentMethod != nil {
		pm := repoModel.PaymentMethod(*order.PaymentMethod)
		result.PaymentMethod = &pm
	}

	if order.TransactionUUID != nil {
		uid := *order.TransactionUUID
		result.TransactionUUID = &uid
	}

	return result
}
