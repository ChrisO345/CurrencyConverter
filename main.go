package main

import (
	"errors"
	"fmt"
	huh "github.com/charmbracelet/huh"
	"log"
	"strconv"
)

var (
	currencyFrom    string
	currencyTo      string
	amount          string
	intAmount       int
	convertedAmount float64
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Currency From:").
				Options(
					huh.NewOption("USD ($)", "usd"),
					huh.NewOption("NZD ($)", "nzd"),
					huh.NewOption("GBP (£)", "gbp"),
					huh.NewOption("EUR (€)", "euro"),
				).
				Value(&currencyFrom),

			huh.NewSelect[string]().
				Title("Currency To:").
				Options(
					huh.NewOption("USD ($)", "usd"),
					huh.NewOption("NZD ($)", "nzd"),
					huh.NewOption("GBP (£)", "gbp"),
					huh.NewOption("EUR (€)", "euro"),
				).
				Value(&currencyTo).
				Validate(func(str string) error {
					if str == currencyFrom {
						return errors.New("currency to cannot be the same as currency from")
					}
					return nil
				}),

			huh.NewInput().
				Title("Amount to convert:").
				Value(&amount).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("amount is required")
					}
					if _, err := strconv.Atoi(str); err != nil {
						return errors.New("amount must be a number")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(currencyFrom)
	fmt.Println(currencyTo)
	fmt.Println(amount)
}
