package entity

import (
	"time"
)

type ProductInput struct {
	Name string `json:"product_name"`
	ProductIn        int8   `json:"product_in"`
	ProductOut       int8   `json:"product_out"`
	User             int8   `json:"user"`

}

type ProductCreated struct {
	ProductName      string `json:"product_name" db:"product_name"`
	ProductIn        int8   `json:"product_in" db:"product_in"`
	ProductOut       int8   `json:"product_out" db:"product_out"`
	Total            int8   `json:"total" db:"total"`
	User             int8   `json:"user" db:"user"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type ProductResponse struct {
	ProductName      string `json:"product_name"`
	ProductIn        int8   `json:"product_in" `
	ProductOut       int8   `json:"product_out"`
	Total            int8   `json:"total" `
	User             int8   `json:"user" `
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" `
	DeletedAt *time.Time `json:"deleted_at" `
}