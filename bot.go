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
	"KGS": true,
}

type Valute struct {
	CharCode string
	Nominal  string
	Value    string
}

type ValCurs struct {
	Valute []Valute
}

var report ValCurs

func main() {

	resp, err := http.Get(CbrUrl)

	if err != nil {
		writeLog(err.Error())
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		writeLog(err.Error())
	}

	regex := regexp.MustCompile("<\\?xml.*?>")

	decoder := charmap.Windows1251.NewDecoder()

	out, decodeErr := decoder.String(string(data))

	if decodeErr != nil {
		writeLog(decodeErr.Error())
	}

	resultData := regex.ReplaceAllString(out, "")

	xml.Unmarshal([]byte(resultData), &report)

	message := ""

	for _, v := range report.Valute {
		if currencies[v.CharCode] {
			message += fmt.Sprintf("%s:\tNominal: %s\t Value: %s\n", v.CharCode, v.Nominal, v.Value)
		}
	}

	sendMessage(message)
}
