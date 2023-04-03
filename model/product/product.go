package otp

import (
	"context"

	"github.com/VincentTC/simple-ecommerce/model"
)

type Repository interface {
	GetAll(ctx context.Context) ([]model.Product, error)
	GetById(ctx context.Context, productId int64) (model.Product, error)
	GetByIds(ctx context.Context, productIds []int64) ([]model.Product, error)
	Insert(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	UpdateQuantity(ctx context.Context, productId int64, newQuantity int64) error
	Delete(ctx context.Context, productId int64) error
}
