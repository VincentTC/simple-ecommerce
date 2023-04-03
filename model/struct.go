package model

import (
	"time"

	"github.com/golang-jwt/jwt"
	"gopkg.in/guregu/null.v3"
)

const (
	OrderPending int = 1
	OrderPaid    int = 2

	RoleCustomer string = "customer"
	RoleAdmin    string = "admin"
)

type Customer struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	IsActive  int       `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt null.Time `db:"updated_at"`
	DeletedAt null.Time `db:"deleted_at"`
}

type Product struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       int64     `db:"price"`
	ImageUrl    string    `db:"image_url"`
	Quantity    int64     `db:"quantity"`
	IsActive    int       `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   null.Time `db:"updated_at"`
	DeletedAt   null.Time `db:"deleted_at"`
}

type Order struct {
	Id         int64     `db:"id"`
	CustomerId int64     `db:"customer_id"`
	TotalPrice int64     `db:"total_price"`
	Status     int       `db:"status"`
	PaidAt     null.Time `db:"paid_at"`
	IsActive   int       `db:"is_active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  null.Time `db:"updated_at"`
	DeletedAt  null.Time `db:"deleted_at"`
}

type OrderProduct struct {
	Id        int64     `db:"id"`
	OrderId   int64     `db:"order_id"`
	ProductId int64     `db:"product_id"`
	Quantity  int64     `db:"quantity"`
	IsActive  int       `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt null.Time `db:"updated_at"`
	DeletedAt null.Time `db:"deleted_at"`
}

//Token struct declaration
type Token struct {
	UserID uint
	Name   string
	Email  string
	*jwt.StandardClaims
}
