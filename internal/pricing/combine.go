package pricing

// CombineModifiers собирает модификаторы в один, применяя их слева направо.
func CombineModifiers(modifiers ...PriceModifier) PriceModifier {
	chain := make([]PriceModifier, len(modifiers))
	copy(chain, modifiers)

	return func(price float64) float64 {
		for _, m := range chain {
			if m == nil {
				continue
			}
			price = m(price)
		}
		return price
	}
}
