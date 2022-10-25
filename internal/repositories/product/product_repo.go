package repositories

import (
	"context"

	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/pkg/postgres"
)

type ProductRepoInterface interface {
	CreateProduct(ctx context.Context, input *entity.ProductCreated) error
	GetProducts(ctx context.Context) ([]entity.ProductResponse, error)
}

type productRepo struct {
	db postgres.Adapter
}

func NewProductRepo(db postgres.Adapter) ProductRepoInterface {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) CreateProduct(ctx context.Context, input *entity.ProductCreated) error {
	query := `INSERT INTO products (product_name, product_in, product_out, total, user_id, created_at, updated_at, deleted_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := r.db.Exec(ctx, query, input.ProductName, input.ProductIn, input.ProductOut, input.Total, input.User, input.CreatedAt, input.UpdatedAt, input.DeletedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) GetProducts(ctx context.Context) ([]entity.ProductResponse, error) {
	query := `SELECT product_name, product_in, product_out, total, user_id, created_at, updated_at, deleted_at  FROM products WHERE deleted_at IS NULL;`

	rows, err := r.db.QueryRows(ctx, query)
	if err != nil {
		return nil, err
	}

	var products []entity.ProductResponse
	for rows.Next() {
		var product entity.ProductResponse
		err := rows.Scan(&product.ProductName, &product.ProductIn, &product.ProductOut, &product.Total, &product.User, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
