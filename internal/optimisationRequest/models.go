package optimisationRequest

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

// Request ...
type Request struct {
	ID           bson.ObjectId  `bson:"_id" json:"id"`
	Name         string         `bson:"name" json:"name"`
	Input        []RequestInput `bson:"input" json:"input"`
	DateCreated  time.Time      `bson:"date_created" json:"date_created"`
}

// RequestInput ...
type RequestInput struct {
	Type     string `bson:"type" json:"type"`
	Location string `bson:"location" json:"location"`
}

// NewRequest ...
type NewRequest struct {
	Name                     string            `json:"name" validate:"required"`
	Input      				 []NewRequestInput `json:"input" validate:"required"`		
}

// NewRequestInput ...
type NewRequestInput struct {
	Type     string `json:"type" validate:"required"`
	Location string `json:"location" validate:"required"`
}