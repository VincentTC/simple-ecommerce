package mysql

import (
	"context"
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

func (m *MySQL) GetAll(ctx context.Context) (res []model.Order, err error) {

	query := `
		SELECT
			id,
			customer_id,
			total_price,
			status,
			paid_at,
			created_at,
			updated_at,
			deleted_at
		FROM 
			orders
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

func (m *MySQL) GetAllPending(ctx context.Context) (res []model.Order, err error) {
	status := model.OrderPending

	query := `
		SELECT
			id,
			customer_id,
			total_price,
			status,
			paid_at,
			created_at,
			updated_at,
			deleted_at
		FROM 
			orders
		WHERE
			status = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, status)

	err = m.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) GetById(ctx context.Context, orderId int64) (res model.Order, err error) {

	query := `
		SELECT
			id,
			customer_id,
			total_price,
			status,
			paid_at,
			created_at,
			updated_at,
			deleted_at
		FROM 
			orders
		WHERE
			id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, orderId)

	err = m.db.GetContext(ctx, &res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) GetByCustomerId(ctx context.Context, customerId int64) (res []model.Order, err error) {

	query := `
		SELECT
			id,
			customer_id,
			total_price,
			status,
			paid_at,
			created_at,
			updated_at,
			deleted_at
		FROM 
			orders
		WHERE
			customer_id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, customerId)

	err = m.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) Insert(ctx context.Context, order *model.Order) error {
	now := time.Now()
	order.CreatedAt = now

	query := `
		INSERT INTO orders (
			customer_id,
			total_price,
			status,
			is_active,
			created_at
		) VALUES (
			?, ?, ?, 1, ?
		)
		;
	`

	res, err := m.db.ExecContext(ctx, query,
		order.CustomerId,
		order.TotalPrice,
		order.Status,
		order.CreatedAt,
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
	order.Id = id

	return nil
}

func (m *MySQL) Update(ctx context.Context, order *model.Order) error {
	now := time.Now()
	order.UpdatedAt.SetValid(now)

	query := `
		UPDATE
			orders
		SET
			customer_id = ?,
			total_price = ?,
			status = ?,
			paid_at = ?,
			updated_at = ?
		WHERE
			id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args,
		order.CustomerId,
		order.TotalPrice,
		order.Status,
		order.PaidAt,
		order.UpdatedAt,
		order.Id,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) Delete(ctx context.Context, orderId int64) error {
	now := time.Now()

	query := `
		UPDATE
			orders
		SET
			is_active = 0,
			deleted_at = ?
		WHERE
			id = ?;
		`

	var args []interface{}
	args = append(args,
		now,
		orderId,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
