package service

import (
	"context"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/repo"
)

type OrderService struct {
	repo repo.IUploadOrder
}

func NewOrderService(repo repo.IUploadOrder) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateUploadedOrder(userID, orderNumber int) (int, error) {
	if !checkOrderNumber(orderNumber) {
		return 0, ErrInvalidNumber
	}

	return s.repo.CreateUploadedOrder(userID, orderNumber)
}

func (s *OrderService) GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrder, error) {
	return s.repo.GetUploadedOrders(ctx, userID)
}

func checkOrderNumber(orderNumber int) bool {
	var luhn int
	number := orderNumber / 10

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}

	return (orderNumber%10+luhn%10)%10 == 0
}
