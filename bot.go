package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"golang.org/x/text/encoding/charmap"
)

const (
	CbrUrl    = "https://cbr.ru/scripts/XML_daily.asp"
	BotApiUrl = "https://api.telegram.org/bot"
)

var currencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"KZT": true,
}

type Valute struct {
	CharCode string
	Nominal  string
	Value    string
}

type ValCurs struct {
	Valute []Valute
}

func main() {

	var report ValCurs

	resp, err := http.Get(CbrUrl)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	regex := regexp.MustCompile("<\\?xml.*?>")

	decoder := charmap.Windows1251.NewDecoder()

	out, decodeErr := decoder.String(string(data))

	if decodeErr != nil {
		fmt.Printf("Decode error: %s", decodeErr)
	}

	resultData := regex.ReplaceAllString(out, "")

	xml.Unmarshal([]byte(resultData), &report)

	for _, v := range report.Valute {
		fmt.Printf("Code: %s\tNominal: %s\t Value: %s\n", v.CharCode, v.Nominal, v.Value)
	}

}
