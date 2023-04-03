package service

import "context"

func New(options *SvOptions, repo Repository) *Service {
	sv := Service{
		options: options,
		repo:    repo,
	}

	return &sv
}

type Sv interface {
	RegisterHandler(ctx context.Context, registerReq RegisterReq) (res RegisterResp, err error)
	LoginHandler(ctx context.Context, loginReq LoginReq) (res LoginResp, err error)
	GetCustomerById(ctx context.Context, customerId int64) (res CustomerResp, err error)

	CreateOrderHandler(ctx context.Context, orderReq OrderReq) (res OrderResp, err error)
	GetAllOrdersHandler(ctx context.Context) (res []OrderResp, err error)
	GetAllPendingOrdersHandler(ctx context.Context) (res []OrderResp, err error)
	GetOrdersByCustomerHandler(ctx context.Context, customerId int64) (res []OrderResp, err error)
	GetOrdersReportHandler(ctx context.Context) (res []OrderResp, err error)
}
