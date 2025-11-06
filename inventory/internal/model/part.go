package model

import "time"

type Part struct {
	Uuid          string            `json:"uuid"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Price         float64           `json:"price"`
	StockQuantity int64             `json:"stock_quantity"`
	Category      Category          `json:"category"`
	Dimensions    *Dimensions       `json:"dimensions"`
	Manufacturer  *Manufacturer     `json:"manufacturer"`
	Tags          []string          `json:"tags"`
	Metadata      map[string]*Value `json:"metadata"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Weight float64 `json:"weight"`
}

type Manufacturer struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Website string `json:"website"`
}

type PartsFilter struct {
	Uuids                 []string   `json:"uuids"`
	Names                 []string   `json:"names"`
	Categories            []Category `json:"categories"`
	ManufacturerCountries []string   `json:"manufacturer_countries"`
	Tags                  []string   `json:"tags"`
}
