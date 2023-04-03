package mysql

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/VincentTC/simple-ecommerce/model"
	"github.com/VincentTC/simple-ecommerce/util/database/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL ...
type MySQL struct {
	db *sqlx.DB
}

// New ...
func New(db *mysql.DB) *MySQL {
	return &MySQL{
		db: db.Client,
	}
}

func (m *MySQL) GetAll(ctx context.Context) (res []model.OrderProduct, err error) {

	query := `
		SELECT
			id,
			order_id,
			product_id,
			quantity,
			created_at,
			updated_at,
			deleted_at
		FROM 
			order_products
		WHERE
			is_active = 1;
		`

	err = m.db.SelectContext(ctx, &res, query)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) GetById(ctx context.Context, orderProductId int64) (res model.OrderProduct, err error) {

	query := `
		SELECT
			id,
			order_id,
			product_id,
			quantity,
			created_at,
			updated_at,
			deleted_at
		FROM 
			order_products
		WHERE
			id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, orderProductId)

	err = m.db.GetContext(ctx, &res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) GetByOrderIds(ctx context.Context, orderIds []int64) (res []model.OrderProduct, err error) {
	if len(orderIds) == 0 {
		return res, errors.New("invalid order ids")
	}

	query := `
		SELECT
			id,
			order_id,
			product_id,
			quantity,
			created_at,
			updated_at,
			deleted_at
		FROM 
			order_products
		WHERE
			order_id IN (?) AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, orderIds)

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	err = m.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) Insert(ctx context.Context, orderProduct *model.OrderProduct) error {
	now := time.Now()
	orderProduct.CreatedAt = now

	query := `
		INSERT INTO order_products (
			order_id,
			product_id,
			quantity,
			is_active,
			created_at
		) VALUES (
			?, ?, ?, 1, ?
		)
		;
	`

	res, err := m.db.ExecContext(ctx, query,
		orderProduct.OrderId,
		orderProduct.ProductId,
		orderProduct.Quantity,
		orderProduct.CreatedAt,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	orderProduct.Id = id

	return nil
}

func (m *MySQL) Update(ctx context.Context, orderProduct *model.OrderProduct) error {
	now := time.Now()
	orderProduct.UpdatedAt.SetValid(now)

	query := `
		UPDATE
			order_products
		SET
			order_id = ?,
			product_id = ?,
			quantity = ?,
			updated_at = ?
		WHERE
			id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args,
		orderProduct.OrderId,
		orderProduct.ProductId,
		orderProduct.Quantity,
		orderProduct.UpdatedAt,
		orderProduct.Id,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) Delete(ctx context.Context, orderProductId int64) error {
	now := time.Now()

	query := `
		UPDATE
			order_products
		SET
			is_active = 0,
			deleted_at = ?
		WHERE
			id = ?;
		`

	var args []interface{}
	args = append(args,
		now,
		orderProductId,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
