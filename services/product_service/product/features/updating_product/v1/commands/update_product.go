package commands

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UpdateProduct struct {
	ProductID   uuid.UUID `validate:"required"`
	Name        string    `validate:"required,gte=0,lte=255"`
	Description string    `validate:"required,gte=0,lte=5000"`
	Price       float64   `validate:"required,gte=0"`
	UpdatedAt   time.Time `validate:"required"`
	Count       int32     `validate:"required,gt=0"`
	InventoryId int64     `validate:"required,gt=0"`
}

func NewUpdateProduct(productID uuid.UUID, name string, description string, price float64, inventoryId int64, count int32) *UpdateProduct {
	return &UpdateProduct{ProductID: productID, Name: name, Description: description,
		Price: price, UpdatedAt: time.Now(), InventoryId: inventoryId, Count: count}
}
