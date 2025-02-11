package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

const (
	CbrUrl    = "https://cbr.ru/scripts/XML_daily.asp"
	BotApiUrl = "https://api.telegram.org/bot"
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 14.2; rv:122.0) Gecko/20100101 Firefox/122.0"
)

var currencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"KGS": true,
}

type Valute struct {
	CharCode  string
	Nominal   string
	Value     string
	VunitRate string
}

type ValCurs struct {
	Valute []Valute
}

var report ValCurs

func main() {

	client := http.Client{}
	req, _ := http.NewRequest("GET", CbrUrl, nil)
	req.Header.Add("User-Agent", UserAgent)

	resp, err := client.Do(req)

	if err != nil {
		writeLog(err.Error())
	}

	data, err := io.ReadAll(resp.Body)

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

			costParts := strings.Split(v.VunitRate, ",")
			costParts[1] = costParts[1][:2]
			finCost := strings.Join(costParts, ".")

			message += fmt.Sprintf("%s:\t %s\n", v.CharCode, finCost)
		}
	}

	sendMessage(message)
}
