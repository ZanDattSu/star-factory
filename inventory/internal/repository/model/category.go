package model

type Category int

const (
	UNKNOWN Category = iota
	ENGINE
	FUEL
	PORTHOLE
	WING
)

func (c Category) String() string {
	return [...]string{"North", "East", "South", "West"}[c]
}
