package repositories

import (
	"context"

	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/pkg/postgres"
)

type RoleRepoInterface interface {
	CreateRole(ctx context.Context, input *entity.RoleCreated) error
}

type roleRepo struct{
	db postgres.Adapter
}

func NewRoleRepo(db postgres.Adapter) RoleRepoInterface{
	return &roleRepo{
		db: db,
	}
}

func (r *roleRepo) CreateRole(ctx context.Context, input *entity.RoleCreated) error{
	query := `INSERT INTO roles (created_at, role_name, updated_at, deleted_at) VALUES($1, $2, $3,$4);`

	_, err := r.db.Exec(ctx, query, input.CreatedAt, input.RoleName, input.UpdatedAt, input.DeletedAt)
	if err != nil{
		return err
	}
	return nil
}