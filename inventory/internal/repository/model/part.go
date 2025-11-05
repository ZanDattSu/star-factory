package model

import "time"

type Part struct {
	Uuid          string            `json:"uuid"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Price         float64           `json:"price"`
	StockQuantity int               `json:"stock_quantity"`
	Category      Category          `json:"category"`
	Dimensions    *Dimensions       `json:"dimensions"`
	Manufacturer  *Manufacturer     `json:"manufacturer"`
	Tags          []string          `json:"tags"`
	Metadata      map[string]string `json:"metadata"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type Dimensions struct {
	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`
	Weight int `json:"weight"`
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
