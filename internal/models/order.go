package models

import "fmt"

// Order — заказ, состоящий из товаров.
type Order struct {
	ID       string
	Customer string
	Items    []*Product
}

// NewOrder создаёт заказ с проверкой полей.
func NewOrder(id, customer string) (*Order, error) {
	if id == "" {
		return nil, fmt.Errorf("models: пустой идентификатор заказа")
	}
	if customer == "" {
		return nil, fmt.Errorf("models: пустой клиент заказа")
	}
	return &Order{ID: id, Customer: customer}, nil
}

// Add добавляет товар и возвращает заказ для цепочки вызовов.
func (o *Order) Add(p *Product) *Order {
	if p != nil {
		o.Items = append(o.Items, p)
	}
	return o
}

// Total возвращает сумму базовых цен товаров.
func (o *Order) Total() float64 {
	var sum float64
	for _, it := range o.Items {
		if it != nil {
			sum += it.Price
		}
	}
	return sum
}

// String реализует fmt.Stringer.
func (o *Order) String() string {
	return fmt.Sprintf("Order{#%s %q items=%d total=%.2f}",
		o.ID, o.Customer, len(o.Items), o.Total())
}
