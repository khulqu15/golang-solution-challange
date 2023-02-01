package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Nebulizer struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Location string             `json:"location,omitempty" validate:"required"`
	Data     []NebulizerData    `json:"data,omitempty"`
}

type NebulizerData struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Nebulizer primitive.ObjectID `json:"nebulizer,omitempty"`
	Power     int                `json:"power,omitempty" validate:"required"`
	Smoke     int                `json:"smoke,omitempty" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty"`
}
