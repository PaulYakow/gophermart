package service

import (
	"errors"
	"github.com/PaulYakow/gophermart/internal/repo"
)

type OrderService struct {
	repo repo.IOrder
}

func NewOrderService(repo repo.IOrder) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(userID, orderNumber int) (int, error) {
	if !checkOrderNumber(orderNumber) {
		return 0, errors.New("invalid order number format")
	}

	return s.repo.CreateOrder(userID, orderNumber)
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
