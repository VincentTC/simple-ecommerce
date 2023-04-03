package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/VincentTC/simple-ecommerce/model"
	mdlCustomer "github.com/VincentTC/simple-ecommerce/model/customer"
	mdlCustomerMysql "github.com/VincentTC/simple-ecommerce/model/customer/mysql"
	mdlOrder "github.com/VincentTC/simple-ecommerce/model/order"
	mdlOrderProduct "github.com/VincentTC/simple-ecommerce/model/order-product"
	mdlOrderProductMysql "github.com/VincentTC/simple-ecommerce/model/order-product/mysql"
	mdlOrderMysql "github.com/VincentTC/simple-ecommerce/model/order/mysql"
	mdlProduct "github.com/VincentTC/simple-ecommerce/model/product"
	mdlProductMysql "github.com/VincentTC/simple-ecommerce/model/product/mysql"
	"github.com/VincentTC/simple-ecommerce/service"
	"github.com/VincentTC/simple-ecommerce/util"
	dbMysql "github.com/VincentTC/simple-ecommerce/util/database/mysql"
)

func scriptReport() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := model.Config{}
	err := util.LoadAndParse("ECOM", &cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
	}
	log.Println("Config successfully loaded")

	//DB
	db, err := dbMysql.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialized database: %v\n", err)
	}
	defer db.Close()
	log.Println("Database successfully initialized")

	appsConfig := cfg.Apps

	// Initializes service
	var customerRepo mdlCustomer.Repository
	customerRepo = mdlCustomerMysql.New(db)

	var productRepo mdlProduct.Repository
	productRepo = mdlProductMysql.New(db)

	var orderRepo mdlOrder.Repository
	orderRepo = mdlOrderMysql.New(db)

	var orderProductRepo mdlOrderProduct.Repository
	orderProductRepo = mdlOrderProductMysql.New(db)

	repo := service.Repository{
		Customer:     customerRepo,
		Product:      productRepo,
		Order:        orderRepo,
		OrderProduct: orderProductRepo,
	}

	svOptions := service.SvOptions{
		AppsConfig: appsConfig,
	}
	var sv service.Sv
	sv = service.New(&svOptions, repo)

	orders, err := sv.GetOrdersReportHandler(context.Background())
	if err != nil {
		log.Fatalf("Failed to get orders data: %v\n", err)
	}

	filenames := "order-report-" + time.Now().Format("2006-01-02")

	log.Println("print to file")
	fs, err := os.OpenFile(filenames+".csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s\n", err)
	}
	defer fs.Close()

	header := "Order ID; Customer Name; Order Date; Total Price of the Order; Status of the Order\n"
	if _, err := fs.WriteString(header); err != nil {
		log.Fatalf("Failed to write file: %s\n", err)
	}

	for _, order := range orders {
		row := fmt.Sprintf("%d; %s; %s; %d; %d\n", order.Id, order.Customer.Name, order.CreatedAt.Format("2006-01-02 15:04:05"), order.TotalPrice, order.Status)
		if _, err := fs.WriteString(row); err != nil {
			log.Fatalf("Failed to write file: %s\n", err)
		}
	}
}
