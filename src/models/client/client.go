package client

import (
	"aherman/src/models/base"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Client todo
type Client struct {
	base.Model
	Name string
	Secret uuid.UUID
}

// Createable represents the fields needed to create a client.
type Createable struct {
	Name string `json:"name"`
}

// Public represents the client fields we'll make publicly available.
type Public struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// BindAttributesFrom TODO
func (public *Public) BindAttributesFrom(client *Client) {
	public.ID = client.ID
	public.Name = client.Name
	public.CreatedAt = client.CreatedAt
	public.UpdatedAt = client.UpdatedAt
	public.DeletedAt = client.DeletedAt
}