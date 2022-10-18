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

func (s *OrderService) GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrderDTO, error) {
	orders, err := s.repo.GetUploadedOrders(ctx, userID)

	result := make([]entity.UploadOrderDTO, len(orders))
	for i, order := range orders {
		result[i] = entity.UploadOrderDTO{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.CreatedAt,
		}
	}

	return result, err
}

func (s *OrderService) CreateWithdrawOrder(userID int, orderNumber string, sum float32) error {
	return s.repo.CreateWithdrawOrder(userID, orderNumber, sum)
}

func (s *OrderService) GetWithdrawOrders(ctx context.Context, userID int) ([]entity.WithdrawOrderDTO, error) {
	orders, err := s.repo.GetWithdrawOrders(ctx, userID)

	result := make([]entity.WithdrawOrderDTO, len(orders))
	for i, order := range orders {
		result[i] = entity.WithdrawOrderDTO{
			Order:       order.Number,
			Sum:         order.Sum,
			ProcessedAt: order.CreatedAt,
		}
	}

	return result, err
}
