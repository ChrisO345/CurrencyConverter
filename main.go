package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
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

	body := apiCall()
	// Free API Key, therefore base is always USD
	data := decodeJson(body)
	generateRatio(data["rates"].(map[string]interface{}))
	fmt.Printf("Converted %.2f %s to %.2f %s\n", intAmount, currencyFrom, convertedAmount, currencyTo)
}

func makeForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Starting Currency: ").
				Options(
					huh.NewOption("USD ($)", "USD"),
					huh.NewOption("NZD ($)", "NZD"),
					huh.NewOption("GBP (£)", "GBP"),
					huh.NewOption("EUR (€)", "EUR"),
				).
				Value(&currencyFrom),

			huh.NewSelect[string]().
				Title("Ending Currency").
				Options(
					huh.NewOption("USD ($)", "USD"),
					huh.NewOption("NZD ($)", "NZD"),
					huh.NewOption("GBP (£)", "GBP"),
					huh.NewOption("EUR (€)", "EUR"),
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

func apiCall() []byte {
	apiKey := os.Getenv("OPENEXCHANGERATES_API_KEY")

	url := "https://openexchangerates.org/api/latest.json?app_id=" + apiKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
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

	return body
}

func decodeJson(body []byte) map[string]interface{} {
	//json.Unmarshal()
	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	return data
}

func generateRatio(rates map[string]interface{}) {
	if currencyFrom == "USD" {
		convertedAmount = intAmount * rates[currencyTo].(float64)
		return
	}
	if currencyTo == "USD" {
		convertedAmount = intAmount / rates[currencyFrom].(float64)
		return
	}
	convertedAmount = intAmount / rates[currencyFrom].(float64) * rates[currencyTo].(float64)
	return
}
