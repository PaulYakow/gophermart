package service

import (
	"context"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/repo"
	"strconv"
)

type OrderService struct {
	repo repo.IOrder
}

func NewOrderService(repo repo.IOrder) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateUploadedOrder(userID, orderNumber int) (int, error) {
	return s.repo.CreateUploadedOrder(userID, strconv.Itoa(orderNumber))
}

func (s *OrderService) GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrder, error) {
	return s.repo.GetUploadedOrders(ctx, userID)
}

func (s *OrderService) GetWithdrawOrders(ctx context.Context, userID int) ([]entity.WithdrawOrder, error) {
	return s.repo.GetWithdrawOrders(ctx, userID)
}
