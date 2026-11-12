package main

import (
	"fmt"

	"go-lab4/internal/pricing"
)

// demoPricing показывает цепочки модификаторов цены.
func demoPricing(base float64) {
	// Пошаговый расчёт: видно, как меняется цена после каждого шага.
	steps := []struct {
		name string
		mod  pricing.PriceModifier
	}{
		{"скидка 10%", pricing.WithDiscount(10)},
		{"налог 20%", pricing.WithTax(20)},
		{"округление до целых", pricing.WithRounding(0)},
		{"минимум 500", pricing.WithMinimumPrice(500)},
	}
	price := base
	rows := [][]string{{"базовая сумма", fmt.Sprintf("%.2f", price), ""}}
	mods := make([]pricing.PriceModifier, 0, len(steps))
	for _, s := range steps {
		before := price
		price = s.mod(price)
		mods = append(mods, s.mod)
		rows = append(rows, []string{s.name, fmt.Sprintf("%.2f", price), fmt.Sprintf("%+.2f", price-before)})
	}
	chain := pricing.CombineModifiers(mods...)
	rows = append(rows, []string{"CombineModifiers (всё сразу)", fmt.Sprintf("%.2f", chain(base)), ""})
	printTable("Цепочка модификаторов по шагам", []string{"Шаг", "Цена", "Изменение"}, rows)

	// Разные цепочки на одной и той же сумме.
	scenarios := []struct {
		name  string
		chain string
		mod   pricing.PriceModifier
	}{
		{"Без скидок", "—", pricing.CombineModifiers()},
		{"Студенческая", "скидка 25, налог 20, округл. 1", pricing.CombineModifiers(
			pricing.WithDiscount(25), pricing.WithTax(20), pricing.WithRounding(1))},
		{"VIP", "скидка 40, налог 20, максимум 400", pricing.CombineModifiers(
			pricing.WithDiscount(40), pricing.WithTax(20), pricing.WithMaximumPrice(400))},
		{"С доставкой", "скидка 10, +150, округл. 0", pricing.CombineModifiers(
			pricing.WithDiscount(10), pricing.WithFixedFee(150), pricing.WithRounding(0))},
		{"Бесплатно", "скидка 100, минимум 50", pricing.CombineModifiers(
			pricing.WithDiscount(100), pricing.WithMinimumPrice(50))},
		{"С nil", "скидка 10, nil, округл. 0", pricing.CombineModifiers(
			pricing.WithDiscount(10), nil, pricing.WithRounding(0))},
	}
	rows = rows[:0]
	for _, s := range scenarios {
		rows = append(rows, []string{s.name, s.chain, fmt.Sprintf("%.2f", s.mod(base))})
	}
	printTable(fmt.Sprintf("Сценарии (база %.2f)", base), []string{"Сценарий", "Цепочка", "Итог"}, rows)

	// Порядок применения важен: надбавка до скидки и после неё дают разный результат.
	a := pricing.CombineModifiers(pricing.WithFixedFee(100), pricing.WithDiscount(50))
	b := pricing.CombineModifiers(pricing.WithDiscount(50), pricing.WithFixedFee(100))
	printTable("Порядок важен (база 1000)", []string{"Порядок", "Итог"}, [][]string{
		{"надбавка 100 -> скидка 50%", fmt.Sprintf("%.2f", a(1000))},
		{"скидка 50% -> надбавка 100", fmt.Sprintf("%.2f", b(1000))},
	})
}

// demoClosures показывает замыкания и захват внешних переменных.
func demoClosures() {
	var rows [][]string
	add := func(what, result string) { rows = append(rows, []string{what, result}) }

	// 1. Замыкание изменяет захваченную переменную.
	count := 0
	increment := func() int {
		count++
		return count
	}
	add("increment() x3", fmt.Sprint(increment(), " ", increment(), " ", increment()))
	add("count снаружи", fmt.Sprint(count))

	// 2. Фабрика: каждый вызов создаёт независимый счётчик со своим состоянием.
	newCounter := func(step int) func() int {
		n := 0
		return func() int {
			n += step
			return n
		}
	}
	c1, c2 := newCounter(1), newCounter(10)
	add("счётчик c1 (шаг 1) x3", fmt.Sprint(c1(), " ", c1(), " ", c1()))
	add("счётчик c2 (шаг 10) x2", fmt.Sprint(c2(), " ", c2()))

	// 3. Анонимная функция, вызываемая сразу.
	add("IIFE func(x){x*x}(7)", fmt.Sprint(func(x int) int { return x * x }(7)))

	// 4. Захват переменной цикла (в Go 1.22+ у каждой итерации своя копия).
	var funcs []func() int
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() int { return i * 10 })
	}
	add("замыкания в цикле i*10", fmt.Sprint(funcs[0](), " ", funcs[1](), " ", funcs[2]()))

	// 5. Кэширующая обёртка над модификатором, со счётчиками вызовов и попаданий.
	calls, hits := 0, 0
	cache := map[float64]float64{}
	withCache := func(fn pricing.PriceModifier) pricing.PriceModifier {
		return func(p float64) float64 {
			calls++
			if v, ok := cache[p]; ok {
				hits++
				return v
			}
			v := fn(p)
			cache[p] = v
			return v
		}
	}
	cached := withCache(pricing.WithDiscount(15))
	add("cached(1000), (1000), (2000)", fmt.Sprint(cached(1000), " ", cached(1000), " ", cached(2000)))
	add("вызовов обёртки / из кэша / записей", fmt.Sprintf("%d / %d / %d", calls, hits, len(cache)))

	// 6. Замыкание как параметризуемый модификатор.
	makeDiscount := func(percent float64) pricing.PriceModifier {
		return func(p float64) float64 { return p * (1 - percent/100) }
	}
	add("makeDiscount(50)(2000)", fmt.Sprintf("%.2f", makeDiscount(50)(2000)))

	// 7. Логирующий декоратор: оборачивает любой модификатор и копит журнал.
	var logRows []string
	withLog := func(name string, fn pricing.PriceModifier) pricing.PriceModifier {
		return func(p float64) float64 {
			r := fn(p)
			logRows = append(logRows, fmt.Sprintf("%s: %.2f -> %.2f", name, p, r))
			return r
		}
	}
	logged := pricing.CombineModifiers(
		withLog("налог 20%", pricing.WithTax(20)),
		withLog("округление", pricing.WithRounding(0)),
	)
	add("логирующий декоратор, итог", fmt.Sprintf("%.2f", logged(99.99)))
	for _, l := range logRows {
		add("  журнал", l)
	}

	printTable("Замыкания", []string{"Пример", "Результат"}, rows)
}
