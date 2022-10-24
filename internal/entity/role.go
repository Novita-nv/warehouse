package entity

import (
	"time"
)

type RoleInput struct {
	Name string `json:"role_name"`
}

type RoleCreated struct {
	RoleName      string `json:"role_name" db:"role_name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}