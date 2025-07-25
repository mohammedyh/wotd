package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	wotdURL   = "https://www.dictionary.com/e/word-of-the-day"
	magenta   = "\033[30;45m"
	underline = "\033[4m"
	noColor   = "\033[0m"
)

func initDocumentReader(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func main() {
	openInBrowser := flag.Bool("o", false, "Opens word of the day page in browser")
	shouldPlayPronunciationAudio := flag.Bool("a", false, "Plays the pronunciation audio for the word")
	flag.Parse()

	if *openInBrowser {
		openWotdPageInBrowser(wotdURL)
	}

	if *shouldPlayPronunciationAudio {
		if err := playPronunciationAudio(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	doc, err := initDocumentReader(wotdURL)
	if err != nil {
		log.Fatalf("error parsing html page: %v", err)
	}

	var (
		word              = doc.Find(".otd-item-headword__word h1.js-fit-text")
		definitionAndType = doc.Find(".otd-item-headword__pos-blocks .otd-item-headword__pos p")
		wordExamples      = doc.Find(".wotd-item-origin__content ul:nth-of-type(2)").Eq(0)
		wordType          = definitionAndType.Eq(0)
		definition        = definitionAndType.Eq(1)
		formattedExamples = strings.TrimSuffix(strings.ReplaceAll(wordExamples.Text(), "\n", " \n - "), "\n - ")
	)

	fmt.Printf("%v: %v\n", colorOutput("Word"), word.First().Text())
	fmt.Printf("%v: %v\n", underlineOutput("Word Type"), strings.Trim(wordType.Text(), " \n"))
	fmt.Printf("%v: %v\n", underlineOutput("Definition"), definition.Text())
	fmt.Printf("%v: %v\n", underlineOutput("Examples"), formattedExamples)
}

func openWotdPageInBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		open("xdg-open", url)
	case "darwin":
		open("open", url)
	case "windows":
		open("start", url)
	}
	os.Exit(0)
}

func open(program, url string) {
	if err := exec.Command(program, url).Run(); err != nil {
		log.Fatal("unable to open word of the day page in browser")
	}
}

func colorOutput(message string) string {
	return fmt.Sprintf("%v %v %v", magenta, message, noColor)
}

func underlineOutput(message string) string {
	return fmt.Sprintf("%v%v%v", underline, message, noColor)
}
