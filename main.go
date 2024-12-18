package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	res, err := http.Get("https://www.dictionary.com/e/word-of-the-day")
	if err != nil {
		log.Fatal("error fetching url:", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("error parsing html page", err)
	}

	word := doc.Find(".otd-item-headword__word h1.js-fit-text")
	definitionAndType := doc.Find(".otd-item-headword__pos-blocks .otd-item-headword__pos p")
	wordExamples := doc.Find(".wotd-item-origin__content ul:nth-of-type(2)").Eq(0)
	wordType := definitionAndType.Eq(0)
	definition := definitionAndType.Eq(1)
	formattedExamples := strings.TrimSuffix(
		strings.ReplaceAll(wordExamples.Text(), "\n", " \n - "), "\n - ",
	)

	fmt.Printf("Word of the day: %v\n", word.First().Text())
	fmt.Printf("Word Type: %s\n", strings.Trim(wordType.Text(), " \n"))
	fmt.Printf("Definition: %v\n", definition.Text())
	fmt.Printf("Examples: %v\n", formattedExamples)
}
