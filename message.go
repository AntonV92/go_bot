package main

import (
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func sendMessage(text string) {
	godotenv.Load()

	token := os.Getenv("TOKEN")
	chatId := os.Getenv("CHAT_ID")

	data := url.Values{
		"chat_id": {chatId},
		"text":    {text},
	}

	messageUrl := BotApiUrl + token + "/sendMessage"
	http.PostForm(messageUrl, data)
}
