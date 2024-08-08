package main

import (
	"errors"
	"fmt"
	huh "github.com/charmbracelet/huh"
	dotenv "github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	currencyFrom    string
	currencyTo      string
	amount          string
	intAmount       float64
	convertedAmount float64
)

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	form := makeForm()

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	intAmount, err = strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(currencyFrom)
	fmt.Println(currencyTo)
	fmt.Println(amount)

	apiCall()
}

func makeForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Starting Currency: ").
				Options(
					huh.NewOption("USD ($)", "usd"),
					huh.NewOption("NZD ($)", "nzd"),
					huh.NewOption("GBP (£)", "gbp"),
					huh.NewOption("EUR (€)", "euro"),
				).
				Value(&currencyFrom),

			huh.NewSelect[string]().
				Title("Ending Currency").
				Options(
					huh.NewOption("USD ($)", "usd"),
					huh.NewOption("NZD ($)", "nzd"),
					huh.NewOption("GBP (£)", "gbp"),
					huh.NewOption("EUR (€)", "euro"),
				).
				Value(&currencyTo).
				Validate(func(str string) error {
					if str == currencyFrom {
						return errors.New("starting and ending currency must be different")
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
					if _, err := strconv.ParseFloat(str, 64); err != nil {
						return errors.New("amount must be a number")
					}
					return nil
				}),
		),
	)
	return form
}

func apiCall() {
	apiKey := os.Getenv("OPENEXCHANGERATES_API_KEY")
	fmt.Println(apiKey)

	url := "https://openexchangerates.org/api/latest.json?app_id=" + apiKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	fmt.Println(string(body))
}
