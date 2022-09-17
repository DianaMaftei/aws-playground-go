package model

import (
	"time"
)

type Order struct {
	Id       int64     `json:"id,omitempty"`
	PetId    int64     `json:"petId,omitempty"`
	Quantity int32     `json:"quantity,omitempty"`
	ShipDate time.Time `json:"shipDate,omitempty"`
	// Order Status
	Status   OrderStatus `json:"status,omitempty"`
	Complete bool        `json:"complete,omitempty"`
}

type OrderStatus string

const (
	Placed    OrderStatus = "placed"
	Approved  OrderStatus = "approved"
	Delivered OrderStatus = "delivered"
)
