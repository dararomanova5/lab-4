package main

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"go-lab4/internal/models"
)

func main() {
	order := demoOrder()
	demoPricing(order.Total())
	demoClosures()
}

// printTable выводит данные в виде простой ASCII-таблицы. Ширина считается в
// символах (рунах), поэтому кириллица выравнивается корректно.
func printTable(title string, headers []string, rows [][]string) {
	fmt.Printf("\n=== %s ===\n", title)

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = utf8.RuneCountInString(h)
	}
	for _, r := range rows {
		for i, c := range r {
			if n := utf8.RuneCountInString(c); n > widths[i] {
				widths[i] = n
			}
		}
	}

	line := func(l, m, r string) {
		parts := make([]string, len(widths))
		for i, w := range widths {
			parts[i] = strings.Repeat("-", w+2)
		}
		fmt.Println(l + strings.Join(parts, m) + r)
	}
	row := func(cells []string) {
		parts := make([]string, len(widths))
		for i, w := range widths {
			c := cells[i]
			pad := strings.Repeat(" ", w-utf8.RuneCountInString(c))
			parts[i] = " " + c + pad + " "
		}
		fmt.Println("|" + strings.Join(parts, "|") + "|")
	}

	line("+", "+", "+")
	row(headers)
	line("+", "+", "+")
	for _, r := range rows {
		row(r)
		line("+", "+", "+")
	}
}

// demoOrder создаёт товары и заказ, показывает работу Stringer.
func demoOrder() *models.Order {
	coffee, err := models.NewProduct("A1", "Кофе", 349)
	if err != nil {
		log.Fatal(err)
	}
	tea, err := models.NewProduct("A2", "Чай", 199.9)
	if err != nil {
		log.Fatal(err)
	}
	cake, err := models.NewProduct("A3", "Пирожное", 150.5)
	if err != nil {
		log.Fatal(err)
	}

	order, err := models.NewOrder("42", "Иванов")
	if err != nil {
		log.Fatal(err)
	}
	order.Add(coffee).Add(tea).Add(cake)

	rows := make([][]string, 0, len(order.Items)+1)
	for _, p := range order.Items {
		rows = append(rows, []string{p.ID, p.Name, fmt.Sprintf("%.2f", p.Price), p.String()})
	}
	rows = append(rows, []string{"", "ИТОГО", fmt.Sprintf("%.2f", order.Total()), order.String()})
	printTable("Заказ", []string{"ID", "Товар", "Цена", "Stringer"}, rows)

	// Проверка валидации: конструкторы должны вернуть ошибку.
	_, err1 := models.NewProduct("X", "Брак", -5)
	_, err2 := models.NewOrder("", "Петров")
	printTable("Валидация", []string{"Вызов", "Ошибка"}, [][]string{
		{`NewProduct("X", "Брак", -5)`, err1.Error()},
		{`NewOrder("", "Петров")`, err2.Error()},
	})
	return order
}
