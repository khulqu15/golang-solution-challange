package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Water struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Location string             `json:"location,omitempty" validate:"required"`
	Data     []WaterData        `json:"data,omitempty"`
}

type WaterData struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Water       primitive.ObjectID `json:"water,omitempty"`
	Acidity     int32              `json:"acidity,omitempty" validate:"required"`
	Oxygen      int32              `json:"oxygen,omitempty" validate:"required"`
	Salt        int32              `json:"salt,omitempty" validate:"required"`
	Temperature int32              `json:"temperature,omitempty" validate:"required"`
	CreatedAt   time.Time          `json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty"`
}
