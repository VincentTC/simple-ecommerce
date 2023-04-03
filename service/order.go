package service

import (
	"context"
	"errors"
	"log"

	"github.com/VincentTC/simple-ecommerce/model"
)

func (s *Service) CreateOrderHandler(ctx context.Context, orderReq OrderReq) (res OrderResp, err error) {
	customerId, ok := ctx.Value(CustomerIdCtxValue).(int64)
	if !ok {
		return res, errors.New("invalid customer id")
	}

	orderReq.CustomerId = customerId

	// get product ids
	var productIds []int64
	productQuantityMap := make(map[int64]int64)
	for _, product := range orderReq.OrderProducts {
		productIds = append(productIds, product.ProductId)
		productQuantityMap[product.ProductId] = product.Quantity
	}

	availableProducts, err := s.repo.Product.GetByIds(ctx, productIds)
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	// validate product ids
	if len(availableProducts) != len(orderReq.OrderProducts) {
		// some product id already inactive
		return res, errors.New("some product id already inactive")
	}

	var totalPrice int64
	productMaps := make(map[int64]*model.Product)

	// check product quantity
	for i, product := range availableProducts {
		if product.Quantity < productQuantityMap[product.Id] {
			// product is out of stock
			return res, errors.New("product is out of stock")
		}

		if productMaps[product.Id] == nil {
			productMaps[product.Id] = &availableProducts[i]
		}

		// calculate total price
		totalPrice += product.Price * productQuantityMap[product.Id]
	}

	orderObj := model.Order{
		CustomerId: orderReq.CustomerId,
		TotalPrice: totalPrice,
		Status:     model.OrderPending,
	}

	// insert order to db
	err = s.repo.Order.Insert(ctx, &orderObj)
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	for _, product := range orderReq.OrderProducts {
		orderProductObj := model.OrderProduct{
			OrderId:   orderObj.Id,
			ProductId: product.ProductId,
			Quantity:  product.Quantity,
		}

		// insert order product to db
		err = s.repo.OrderProduct.Insert(ctx, &orderProductObj)
		if err != nil {
			log.Println(err.Error())
			return res, err
		}

		newQuantity := productMaps[product.ProductId].Quantity - product.Quantity

		// update product quantity
		err = s.repo.Product.UpdateQuantity(ctx, product.ProductId, newQuantity)
		if err != nil {
			log.Println(err.Error())
			return res, err
		}
	}

	orderResp, err := s.repo.Order.GetById(ctx, orderObj.Id)
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	var orderProducts []OrderProduct

	for i := range orderReq.OrderProducts {
		orderProducts = append(orderProducts, OrderProduct{
			Product: Product{
				Id:          productMaps[orderReq.OrderProducts[i].ProductId].Id,
				Name:        productMaps[orderReq.OrderProducts[i].ProductId].Name,
				Description: productMaps[orderReq.OrderProducts[i].ProductId].Description,
				Price:       productMaps[orderReq.OrderProducts[i].ProductId].Price,
				ImageUrl:    productMaps[orderReq.OrderProducts[i].ProductId].ImageUrl,
			},
			Quantity: orderReq.OrderProducts[i].Quantity,
		})
	}

	res = OrderResp{
		Id: orderResp.Id,
		Customer: CustomerResp{
			Id: orderResp.CustomerId,
		},
		TotalPrice:    orderResp.TotalPrice,
		OrderProducts: orderProducts,
		Status:        orderResp.Status,
		PaidAt:        orderResp.PaidAt,
		CreatedAt:     orderResp.CreatedAt,
	}

	return res, nil
}

func (s *Service) GetAllOrdersHandler(ctx context.Context) (res []OrderResp, err error) {

	// get all orders
	orders, err := s.repo.Order.GetAll(ctx)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	res, err = s.generateOrderProducts(ctx, orders)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	return res, nil
}

func (s *Service) GetAllPendingOrdersHandler(ctx context.Context) (res []OrderResp, err error) {

	// get all orders
	orders, err := s.repo.Order.GetAllPending(ctx)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	res, err = s.generateOrderProducts(ctx, orders)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	return res, nil
}

func (s *Service) GetOrdersByCustomerHandler(ctx context.Context, customerId int64) (res []OrderResp, err error) {
	customerIdCtx, ok := ctx.Value(CustomerIdCtxValue).(int64)
	if !ok {
		return res, errors.New("invalid customer id")
	}

	if customerId != customerIdCtx {
		return res, errors.New("invalid customer id")
	}

	// get orders by customer id
	orders, err := s.repo.Order.GetByCustomerId(ctx, customerId)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	res, err = s.generateOrderProducts(ctx, orders)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	return res, nil
}

func (s *Service) generateOrderProducts(ctx context.Context, orders []model.Order) (res []OrderResp, err error) {

	var orderIds []int64
	var customerIds []int64
	customerIdsMap := make(map[int64]bool)
	for _, order := range orders {
		orderIds = append(orderIds, order.Id)

		if !customerIdsMap[order.CustomerId] {
			customerIds = append(customerIds, order.CustomerId)
			customerIdsMap[order.CustomerId] = true
		}
	}

	orderProducts, err := s.repo.OrderProduct.GetByOrderIds(ctx, orderIds)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	// get products data
	var productIds []int64
	productIdsMap := make(map[int64]bool)
	orderProductsMap := make(map[int64][]OrderProduct)
	for _, orderProduct := range orderProducts {
		if !productIdsMap[orderProduct.ProductId] {
			productIds = append(productIds, orderProduct.ProductId)
			productIdsMap[orderProduct.ProductId] = true
		}

		orderProductsMap[orderProduct.OrderId] = append(orderProductsMap[orderProduct.OrderId], OrderProduct{
			Product: Product{
				Id: orderProduct.ProductId,
			},
			Quantity: orderProduct.Quantity,
		})
	}

	products, err := s.repo.Product.GetByIds(ctx, productIds)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}
	productMaps := make(map[int64]*model.Product)
	for i, product := range products {
		if productMaps[product.Id] == nil {
			productMaps[product.Id] = &products[i]
		}
	}

	customers, err := s.repo.Customer.GetByIds(ctx, customerIds)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}
	customerMaps := make(map[int64]*model.Customer)
	for i, customer := range customers {
		if customerMaps[customer.Id] == nil {
			customerMaps[customer.Id] = &customers[i]
		}
	}

	for _, order := range orders {
		orderProducts := orderProductsMap[order.Id]

		for i := range orderProducts {
			orderProducts[i].Product = Product{
				Id:          productMaps[orderProducts[i].Product.Id].Id,
				Name:        productMaps[orderProducts[i].Product.Id].Name,
				Description: productMaps[orderProducts[i].Product.Id].Description,
				Price:       productMaps[orderProducts[i].Product.Id].Price,
				ImageUrl:    productMaps[orderProducts[i].Product.Id].ImageUrl,
			}
		}

		orderRes := OrderResp{
			Id: order.Id,
			Customer: CustomerResp{
				Id:    order.CustomerId,
				Name:  customerMaps[order.CustomerId].Name,
				Email: customerMaps[order.CustomerId].Email,
			},
			TotalPrice:    order.TotalPrice,
			OrderProducts: orderProducts,
			Status:        order.Status,
			PaidAt:        order.PaidAt,
			CreatedAt:     order.CreatedAt,
		}

		res = append(res, orderRes)
	}

	return res, nil
}

func (s *Service) GetOrdersReportHandler(ctx context.Context) (res []OrderResp, err error) {

	// get all orders
	orders, err := s.repo.Order.GetAll(ctx)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	res, err = s.generateOrderProducts(ctx, orders)
	if err != nil {
		log.Println(err.Error())
		return res, nil
	}

	return res, nil
}
