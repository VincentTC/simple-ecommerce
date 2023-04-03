package main

import (
	"log"
	"time"

	"github.com/VincentTC/simple-ecommerce/controller"
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
	"github.com/VincentTC/simple-ecommerce/util/router"
	"github.com/VincentTC/simple-ecommerce/util/webserver"
)

func server() {
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

	webOptions := &webserver.Options{
		ListenAddress:  cfg.Server.HTTPPort,
		MaxConnections: 1024,
		ReadTimeout:    30 * time.Second,
	}

	// Initializes web handler
	webServer := webserver.NewWithEngine(router.Muxie, webOptions)
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

	// Controller options
	ctlOptions := &controller.Options{
		Service:    sv,
		AppsConfig: appsConfig,
	}
	ctl := controller.New(ctlOptions)
	ctl.Register(*webServer.Router())

	log.Printf("Running application at port %s...", cfg.Server.HTTPPort)

	if err := webServer.RunGraceful(); err != nil {
		panic(err)
	}
}
