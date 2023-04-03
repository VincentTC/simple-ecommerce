package main

import (
	"log"

	_cron "github.com/VincentTC/simple-ecommerce/cron"
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

func cron() {
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

	// Initialize Cron
	cron := _cron.NewCron(&_cron.Options{
		Sv: sv,
		CronSchedule: _cron.CronSchedule{
			OrderReminderSchedule: cfg.CronSchedule.OrderReminder,
		},
		Apps: appsConfig,
	})
	defer cron.Stop()

	cron.Init()
	log.Println("Cron successfully initialized")

	log.Println("Running Cron...")
	cron.Run()
}
