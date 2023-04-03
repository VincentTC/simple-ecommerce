package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/VincentTC/simple-ecommerce/model"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterHandler(ctx context.Context, registerReq RegisterReq) (res RegisterResp, err error) {

	_, err = s.repo.Customer.GetByEmail(ctx, registerReq.Email)
	if err == nil {
		err = errors.New("email already registered")
		return res, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return res, err
	}

	customerObj := model.Customer{
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: string(password),
		Role:     model.RoleCustomer,
	}

	// insert customer to db
	err = s.repo.Customer.Insert(ctx, &customerObj)
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	res = RegisterResp{
		Success: true,
	}

	return res, nil
}

func (s *Service) LoginHandler(ctx context.Context, loginReq LoginReq) (res LoginResp, err error) {

	customer, err := s.repo.Customer.GetByEmail(ctx, loginReq.Email)
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(loginReq.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		// password does not match
		log.Println(err.Error())
		return res, err
	}

	expiresAt := time.Now().Add(time.Hour * 1).Unix()

	tk := &jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", customer.Id),
		ExpiresAt: expiresAt,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(s.options.AppsConfig.AuthSecret))
	if err != nil {
		log.Println(err)
		return res, err
	}

	res = LoginResp{
		Success:     true,
		AccessToken: tokenString,
	}

	return res, nil
}

func (s *Service) GetCustomerById(ctx context.Context, customerId int64) (res CustomerResp, err error) {

	customer, err := s.repo.Customer.GetById(ctx, customerId)
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	res = CustomerResp{
		Id:    customer.Id,
		Name:  customer.Name,
		Email: customer.Email,
		Role:  customer.Role,
	}

	return res, nil
}
