package otp

import (
	"context"

	"github.com/VincentTC/simple-ecommerce/model"
)

type Repository interface {
	GetAll(ctx context.Context) ([]model.Customer, error)
	GetById(ctx context.Context, customerId int64) (model.Customer, error)
	GetByIds(ctx context.Context, customerIds []int64) ([]model.Customer, error)
	GetByEmail(ctx context.Context, email string) (model.Customer, error)
	Insert(ctx context.Context, customer *model.Customer) error
	Update(ctx context.Context, customer *model.Customer) error
	Delete(ctx context.Context, customerId int64) error
}
