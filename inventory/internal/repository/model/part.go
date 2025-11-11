package model

import (
	"time"

	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

// Part - модель детали в MongoDB
type Part struct {
	Uuid          string            `json:"uuid" bson:"uuid"`
	Name          string            `json:"name" bson:"name"`
	Description   string            `json:"description" bson:"description"`
	Price         float64           `json:"price" bson:"price"`
	StockQuantity int64             `json:"stock_quantity" bson:"stock_quantity"`
	Category      Category          `json:"category" bson:"category"`
	Dimensions    *Dimensions       `json:"dimensions" bson:"dimensions, omitempty"`
	Manufacturer  *Manufacturer     `json:"manufacturer" bson:"manufacturer, omitempty"`
	Tags          []string          `json:"tags" bson:"tags"`
	Metadata      map[string]*Value `json:"metadata" bson:"metadata, omitempty"`
	CreatedAt     time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" bson:"updated_at"`
}

type Dimensions struct {
	Length float64 `json:"length" bson:"length"`
	Width  float64 `json:"width" bson:"width"`
	Height float64 `json:"height" bson:"height"`
	Weight float64 `json:"weight" bson:"weight"`
}

type Manufacturer struct {
	Name    string `json:"name" bson:"name"`
	Country string `json:"country" bson:"country"`
	Website string `json:"website" bson:"website"`
}

type Category string

const (
	CategoryUnspecified Category = "UNSPECIFIED"
	CategoryEngine      Category = "ENGINE"
	CategoryFuel        Category = "FUEL"
	CategoryPorthole    Category = "PORTHOLE"
	CategoryWing        Category = "WING"
)

type PartsFilter struct {
	Uuids                 []string   `json:"uuids" bson:"uuids"`
	Names                 []string   `json:"names" bson:"names"`
	Categories            []Category `json:"categories" bson:"categories"`
	ManufacturerCountries []string   `json:"manufacturer_countries" bson:"manufacturer_countries"`
	Tags                  []string   `json:"tags" bson:"tags"`
}
