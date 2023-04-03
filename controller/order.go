package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/VincentTC/simple-ecommerce/service"
	"github.com/VincentTC/simple-ecommerce/util/router"
)

func (ctl *Ctl) createOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		code := http.StatusBadRequest
		http.Error(w, "Form not valid!", code)
		return
	}
	defer r.Body.Close()

	var orderReq service.OrderReq
	if err := json.Unmarshal(body, &orderReq); err != nil {
		code := http.StatusBadRequest
		http.Error(w, "Form not valid!", code)
		return
	}

	respObj, err := ctl.sv.CreateOrderHandler(ctx, orderReq)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}

	resp, err := json.Marshal(respObj)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (ctl *Ctl) getAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	respObj, err := ctl.sv.GetAllOrdersHandler(ctx)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}

	resp, err := json.Marshal(respObj)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (ctl *Ctl) getOrdersByCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	customerId, err := strconv.ParseInt(router.GetParam(r, "customerId"), 10, 64)
	if err != nil {
		code := http.StatusBadRequest
		http.Error(w, "Customer Id not valid!", code)
		return
	}

	respObj, err := ctl.sv.GetOrdersByCustomerHandler(ctx, customerId)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}

	resp, err := json.Marshal(respObj)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
