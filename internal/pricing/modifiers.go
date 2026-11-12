// Package pricing содержит модификаторы цены и их композицию.
package pricing

import "math"

// PriceModifier — любое преобразование цены.
type PriceModifier func(float64) float64

// WithDiscount уменьшает цену на percent %.
func WithDiscount(percent float64) PriceModifier {
	switch {
	case percent < 0:
		percent = 0
	case percent > 100:
		percent = 100
	}
	mult := 1 - percent/100
	return func(price float64) float64 {
		return price * mult
	}
}

// WithTax увеличивает цену на rate %.
func WithTax(rate float64) PriceModifier {
	if rate < 0 {
		rate = 0
	}
	mult := 1 + rate/100
	return func(price float64) float64 {
		return price * mult
	}
}

// WithRounding округляет цену до precision знаков после запятой.
func WithRounding(precision int) PriceModifier {
	coef := math.Pow10(precision)
	return func(price float64) float64 {
		return math.Round(price*coef) / coef
	}
}

// WithMinimumPrice поднимает цену до minPrice, если она ниже.
func WithMinimumPrice(minPrice float64) PriceModifier {
	return func(price float64) float64 {
		if price < minPrice {
			return minPrice
		}
		return price
	}
}

// WithMaximumPrice ограничивает цену сверху.
func WithMaximumPrice(maxPrice float64) PriceModifier {
	return func(price float64) float64 {
		if price > maxPrice {
			return maxPrice
		}
		return price
	}
}

// WithFixedFee добавляет фиксированную надбавку.
func WithFixedFee(fee float64) PriceModifier {
	return func(price float64) float64 {
		return price + fee
	}
}
