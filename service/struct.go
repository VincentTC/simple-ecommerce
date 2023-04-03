package service

import (
	"time"

	"github.com/VincentTC/simple-ecommerce/model"
	mdlCustomer "github.com/VincentTC/simple-ecommerce/model/customer"
	mdlOrder "github.com/VincentTC/simple-ecommerce/model/order"
	mdlOrderProduct "github.com/VincentTC/simple-ecommerce/model/order-product"
	mdlProduct "github.com/VincentTC/simple-ecommerce/model/product"
	"gopkg.in/guregu/null.v3"
)

type CustomerIdCtx string

const (
	CustomerIdCtxValue CustomerIdCtx = "customerId"
)

type SvOptions struct {
	AppsConfig model.AppsConfig
}

type Service struct {
	options *SvOptions
	repo    Repository
}

type Repository struct {
	Customer     mdlCustomer.Repository
	Product      mdlProduct.Repository
	Order        mdlOrder.Repository
	OrderProduct mdlOrderProduct.Repository
}

type (
	RegisterReq struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	RegisterResp struct {
		Success bool `json:"success"`
	}

	LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResp struct {
		Success     bool   `json:"success"`
		AccessToken string `json:"access_token"`
	}

	CustomerResp struct {
		Id    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	OrderReq struct {
		CustomerId    int64             `json:"customer_id"`
		OrderProducts []OrderProductReq `json:"order_products"`
	}

	OrderProductReq struct {
		ProductId int64 `json:"product_id"`
		Quantity  int64 `json:"quantity"`
	}

	OrderResp struct {
		Id            int64          `json:"id"`
		Customer      CustomerResp   `json:"customer"`
		TotalPrice    int64          `json:"total_price"`
		OrderProducts []OrderProduct `json:"order_products"`
		Status        int            `json:"status"`
		PaidAt        null.Time      `json:"paid_at"`
		CreatedAt     time.Time      `json:"created_at"`
	}

	OrderProduct struct {
		Product  Product `json:"product"`
		Quantity int64   `json:"quantity"`
	}

	Product struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int64  `json:"price"`
		ImageUrl    string `json:"image_url"`
	}
)
