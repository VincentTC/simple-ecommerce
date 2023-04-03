package otp

import (
	"context"

	"github.com/VincentTC/simple-ecommerce/model"
)

type Repository interface {
	GetAll(ctx context.Context) ([]model.OrderProduct, error)
	GetById(ctx context.Context, orderProductId int64) (model.OrderProduct, error)
	GetByOrderIds(ctx context.Context, orderIds []int64) ([]model.OrderProduct, error)
	Insert(ctx context.Context, orderProduct *model.OrderProduct) error
	Update(ctx context.Context, orderProduct *model.OrderProduct) error
	Delete(ctx context.Context, orderProductId int64) error
}
