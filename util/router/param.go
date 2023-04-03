package router

import (
	"net/http"

	"github.com/VincentTC/simple-ecommerce/util/router/internal/param"
)

// GetParam returns param k value
func GetParam(r *http.Request, k string) string {
	return param.GetParam(r, k)
}
