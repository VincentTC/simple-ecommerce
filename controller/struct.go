package controller

import (
	"github.com/VincentTC/simple-ecommerce/model"
	"github.com/VincentTC/simple-ecommerce/service"
)

type Options struct {
	Service    service.Sv
	AppsConfig model.AppsConfig
}

type Ctl struct {
	sv service.Sv

	AppsConfig model.AppsConfig
}
