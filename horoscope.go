package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Horoscope struct {
	Sign    string    `json:"sign"`
	Content []Section `json:"sections"`
}

const URL = "https://www.20minutes.fr/horoscope/horoscope-"

func retrieveSign(sign string) Horoscope {
	// Retrieves the horoscope for a specific sign
	var output Horoscope
	var allSections []Section

	collector := colly.NewCollector()
	urlToRequest := URL + sign

	collector.OnHTML("div.c-horoscope-topic-card", func(e *colly.HTMLElement) {
		section := Section{}
		section.Title = e.ChildText("h3")
		section.Content = e.ChildText("p")

		allSections = append(allSections, section)
	})

	collector.OnScraped(func(r *colly.Response) {
		horoscope := Horoscope{Sign: sign, Content: allSections}
		output = horoscope
	})

	err := collector.Visit(urlToRequest)
	if err != nil {
		log.Println("Unable to request", urlToRequest)
	}
	return output
}

func main() {
	signs := [12]string{"belier", "scorpion", "lion", "cancer", "capricorne", "taureau", "poisson", "balance", "sagittaire", "vierge", "verseau", "gemeaux"}

	allHoroscopes := make([]Horoscope, 0)

	for _, s := range signs {
		horoscope := retrieveSign(s)
		allHoroscopes = append(allHoroscopes, horoscope)
	}

	WriteJSON(allHoroscopes)
}

func WriteJSON(data []Horoscope) {
	currentDate := time.Now().Format("20060101")
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}
	_ = os.WriteFile("horoscope_"+currentDate+".json", file, 0644)
}
