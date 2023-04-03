package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VincentTC/simple-ecommerce/service"
	"github.com/VincentTC/simple-ecommerce/util/limiter"
	"github.com/VincentTC/simple-ecommerce/util/router"
	"github.com/golang-jwt/jwt"
)

func New(options *Options) *Ctl {
	return &Ctl{
		sv:         options.Service,
		AppsConfig: options.AppsConfig,
	}
}

func (ctl *Ctl) Register(router router.Router) {
	// limit 100 requests by ip for 1 minute
	r := router.Use(limiter.LimitByIP(100, 1*time.Minute))

	r.POST("/v1/register", ctl.register)
	r.POST("/v1/login", ctl.login)

	r.GET("/v1/orders", ctl.verifyAuth(ctl.getAllOrders, "admin"))
	r.GET("/v1/orders/customer/:customerId", ctl.verifyAuth(ctl.getOrdersByCustomer, "customer"))
	r.POST("/v1/orders", ctl.verifyAuth(ctl.createOrder, "customer"))
}

func (ctl *Ctl) verifyAuth(next http.HandlerFunc, role string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqToken := r.Header.Get("Authorization")
		if reqToken == "" {
			code := http.StatusUnauthorized
			http.Error(w, "Token not found!", code)
			return
		}
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) <= 0 {
			code := http.StatusUnauthorized
			http.Error(w, "Token not found!", code)
			return
		}
		reqToken = splitToken[1]

		tk := &jwt.StandardClaims{}

		token, err := jwt.ParseWithClaims(reqToken, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(ctl.AppsConfig.AuthSecret), nil
		})
		if err != nil {
			log.Println(err)
			code := http.StatusUnauthorized
			http.Error(w, err.Error(), code)
			return
		}

		payload := token.Claims.(*jwt.StandardClaims)
		customerId, err := strconv.ParseInt(payload.Subject, 10, 64)
		if err != nil {
			log.Println(err)
			code := http.StatusUnauthorized
			http.Error(w, err.Error(), code)
			return
		}

		customer, err := ctl.sv.GetCustomerById(r.Context(), customerId)
		if err != nil {
			log.Println(err)
			code := http.StatusUnauthorized
			http.Error(w, err.Error(), code)
			return
		}
		if customer.Role != role {
			code := http.StatusUnauthorized
			http.Error(w, "Unauthorized", code)
			return
		}

		ctx := context.WithValue(r.Context(), service.CustomerIdCtxValue, customerId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
