package models

import "fmt"

// Product — товар в каталоге.
type Product struct {
	ID    string
	Name  string
	Price float64
}

// NewProduct создаёт товар с проверкой полей.
func NewProduct(id, name string, price float64) (*Product, error) {
	switch {
	case id == "":
		return nil, fmt.Errorf("models: пустой идентификатор товара")
	case name == "":
		return nil, fmt.Errorf("models: пустое название товара")
	case price < 0:
		return nil, fmt.Errorf("models: отрицательная цена %.2f", price)
	}
	return &Product{ID: id, Name: name, Price: price}, nil
}

// String реализует fmt.Stringer.
func (p *Product) String() string {
	return fmt.Sprintf("Product{#%s %q %.2f}", p.ID, p.Name, p.Price)
}
