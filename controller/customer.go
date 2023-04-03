package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/VincentTC/simple-ecommerce/service"
)

func (ctl *Ctl) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		code := http.StatusBadRequest
		http.Error(w, "Form not valid!", code)
		return
	}
	defer r.Body.Close()

	var registerReq service.RegisterReq
	if err := json.Unmarshal(body, &registerReq); err != nil {
		code := http.StatusBadRequest
		http.Error(w, "Form not valid!", code)
		return
	}

	respObj, err := ctl.sv.RegisterHandler(ctx, registerReq)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	resp, err := json.Marshal(respObj)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (ctl *Ctl) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		code := http.StatusBadRequest
		http.Error(w, "Form not valid!", code)
		return
	}
	defer r.Body.Close()

	var loginReq service.LoginReq
	if err := json.Unmarshal(body, &loginReq); err != nil {
		code := http.StatusBadRequest
		http.Error(w, "Form not valid!", code)
		return
	}

	respObj, err := ctl.sv.LoginHandler(ctx, loginReq)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	resp, err := json.Marshal(respObj)
	if err != nil {
		log.Println(err.Error())
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
