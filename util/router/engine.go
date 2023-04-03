package router

import (
	"github.com/VincentTC/simple-ecommerce/util/router/internal/httprouter"
	"github.com/VincentTC/simple-ecommerce/util/router/internal/muxie"
)

// EngineType ...
type EngineType uint8

const (
	// HTTPRouter ...
	HTTPRouter EngineType = iota + 1
	// Muxie ...
	Muxie
)

func getRouterEngine(et EngineType) engine {
	switch et {
	case HTTPRouter:
		return httprouter.New()
	case Muxie:
		return muxie.New()
	default:
		panic("Router engine is not found")
	}
}
