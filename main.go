package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	MAGENTA   = "\033[30;45m"
	UNDERLINE = "\033[4m"
	NOCOLOR   = "\033[0m"
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

	fmt.Printf("%v: %v\n", colorOutput("Word"), word.First().Text())
	fmt.Printf("%v: %v\n", underlineOutput("Word Type"), strings.Trim(wordType.Text(), " \n"))
	fmt.Printf("%v: %v\n", underlineOutput("Definition"), definition.Text())
	fmt.Printf("%v: %v\n", underlineOutput("Examples"), formattedExamples)
}

func colorOutput(message string) string {
	return fmt.Sprintf("%v %v %v", MAGENTA, message, NOCOLOR)
}

func underlineOutput(message string) string {
	return fmt.Sprintf("%v%v%v", UNDERLINE, message, NOCOLOR)
}
