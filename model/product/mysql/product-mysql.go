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

func (m *MySQL) GetAll(ctx context.Context) (res []model.Product, err error) {

	query := `
		SELECT
			id,
			name,
			description,
			price,
			image_url,
			quantity,
			created_at,
			updated_at,
			deleted_at
		FROM 
			products
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

func (m *MySQL) GetById(ctx context.Context, productId int64) (res model.Product, err error) {

	query := `
		SELECT
			id,
			name,
			description,
			price,
			image_url,
			quantity,
			created_at,
			updated_at,
			deleted_at
		FROM 
			products
		WHERE
			id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, productId)

	err = m.db.GetContext(ctx, &res, query, args...)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

func (m *MySQL) GetByIds(ctx context.Context, productIds []int64) (res []model.Product, err error) {
	if len(productIds) == 0 {
		return res, errors.New("invalid product ids")
	}

	query := `
		SELECT
			id,
			name,
			description,
			price,
			image_url,
			quantity,
			created_at,
			updated_at,
			deleted_at
		FROM 
			products
		WHERE
			id IN (?) AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args, productIds)

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

func (m *MySQL) Insert(ctx context.Context, product *model.Product) error {
	now := time.Now()
	product.CreatedAt = now

	query := `
		INSERT INTO products (
			name,
			description,
			price,
			image_url,
			quantity,
			is_active,
			created_at
		) VALUES (
			?, ?, ?, ?, ?, 1, ?
		)
		;
	`

	res, err := m.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.ImageUrl,
		product.Quantity,
		product.CreatedAt,
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
	product.Id = id

	return nil
}

func (m *MySQL) Update(ctx context.Context, product *model.Product) error {
	now := time.Now()
	product.UpdatedAt.SetValid(now)

	query := `
		UPDATE
			products
		SET
			name = ?,
			description = ?,
			price = ?,
			image_url = ?,
			quantity = ?,
			updated_at = ?
		WHERE
			id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args,
		product.Name,
		product.Description,
		product.Price,
		product.ImageUrl,
		product.Quantity,
		product.UpdatedAt,
		product.Id,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) UpdateQuantity(ctx context.Context, productId, quantity int64) error {
	now := time.Now()

	query := `
		UPDATE
			products
		SET
			quantity = ?,
			updated_at = ?
		WHERE
			id = ? AND
			is_active = 1;
		`

	var args []interface{}
	args = append(args,
		quantity,
		now,
		productId,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) Delete(ctx context.Context, productId int64) error {
	now := time.Now()

	query := `
		UPDATE
			products
		SET
			is_active = 0,
			deleted_at = ?
		WHERE
			id = ?;
		`

	var args []interface{}
	args = append(args,
		now,
		productId,
	)

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
