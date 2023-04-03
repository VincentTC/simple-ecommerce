package otp

import (
	"context"

	"github.com/VincentTC/simple-ecommerce/model"
)

type Repository interface {
	GetAll(ctx context.Context) ([]model.Order, error)
	GetAllPending(ctx context.Context) ([]model.Order, error)
	GetById(ctx context.Context, orderId int64) (model.Order, error)
	GetByCustomerId(ctx context.Context, customerId int64) ([]model.Order, error)
	Insert(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, orderId int64) error
}
