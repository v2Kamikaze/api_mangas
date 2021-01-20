package main

import (
	"api_refeita/scraper"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mh := scraper.NewScraper("https://mangahosted.com/mangas/page/", http.Client{Timeout: time.Second * 60})
	titles := mh.GetAllTitles(229)
	file, err := os.Create("../data/titles.json")
	if err != nil {
		log.Fatal("Erro ao criar arquivo! Erro: ", err.Error())
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(titles)
	if err != nil {
		log.Fatal("Erro ao escrever no arquivo. Erro: ", err.Error())
	}
}
