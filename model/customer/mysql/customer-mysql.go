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

func (m *MySQL) GetAll(ctx context.Context) (res []model.Customer, err error) {

	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			created_at,
			updated_at,
			deleted_at
		FROM 
			customers
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

func (m *MySQL) GetById(ctx context.Context, customerId int64) (res model.Customer, err error) {

	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			created_at,
			updated_at,
			deleted_at
		FROM 
			customers
		WHERE
			id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, customerId)

	err = m.db.GetContext(ctx, &res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) GetByIds(ctx context.Context, customerIds []int64) (res []model.Customer, err error) {
	if len(customerIds) == 0 {
		return res, errors.New("invalid customer ids")
	}

	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			created_at,
			updated_at,
			deleted_at
		FROM 
			customers
		WHERE
			id IN (?) AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, customerIds)

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

func (m *MySQL) GetByEmail(ctx context.Context, email string) (res model.Customer, err error) {

	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			created_at,
			updated_at,
			deleted_at
		FROM 
			customers
		WHERE
			email = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, email)

	err = m.db.GetContext(ctx, &res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) Insert(ctx context.Context, customer *model.Customer) error {
	now := time.Now()
	customer.CreatedAt = now

	query := `
		INSERT INTO customers (
			name,
			email,
			password,
			role,
			is_active,
			created_at
		) VALUES (
			?, ?, ?, ?, 1, ?
		)
		;
	`

	res, err := m.db.ExecContext(ctx, query,
		customer.Name,
		customer.Email,
		customer.Password,
		customer.Role,
		customer.CreatedAt,
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
	customer.Id = id

	return nil
}

func (m *MySQL) Update(ctx context.Context, customer *model.Customer) error {
	now := time.Now()
	customer.UpdatedAt.SetValid(now)

	query := `
		UPDATE
			customers
		SET
			name = ?,
			email = ?,
			password = ?,
			role = ?,
			updated_at = ?
		WHERE
			id = ?;
		`

	var args []interface{}
	args = append(args,
		customer.Name,
		customer.Email,
		customer.Password,
		customer.Role,
		customer.UpdatedAt,
		customer.Id,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) Delete(ctx context.Context, customerId int64) error {
	now := time.Now()

	query := `
		UPDATE
			customers
		SET
			is_active = 0,
			deleted_at = ?
		WHERE
			id = ?;
		`

	var args []interface{}
	args = append(args,
		now,
		customerId,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
