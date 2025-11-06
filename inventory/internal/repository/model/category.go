package model

type Category int

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

func (c Category) String() string {
	return [...]string{
		"Unknown",
		"Engine",
		"Fuel",
		"Porthole",
		"Wing",
	}[c]
}
