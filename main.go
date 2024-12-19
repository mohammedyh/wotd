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
	WOTD_URL  = "https://www.dictionary.com/e/word-of-the-day"
	MAGENTA   = "\033[30;45m"
	UNDERLINE = "\033[4m"
	NOCOLOR   = "\033[0m"
)

func main() {
	openInBrowserFlag := flag.Bool("o", false, "Opens word of the day page in browser")

	flag.Parse()

	if *openInBrowserFlag {
		openInBrowser(WOTD_URL)
		os.Exit(0)
	}

	res, err := http.Get(WOTD_URL)
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

func openInBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		open("xdg-open", url)
	case "darwin":
		open("open", url)
	case "windows":
		open("start", url)
	}
}

func open(program, url string) {
	cmd := exec.Command(program, url)
	err := cmd.Run()
	if err != nil {
		log.Fatal("unable to open word of the day page in browser")
	}
}

func colorOutput(message string) string {
	return fmt.Sprintf("%v %v %v", MAGENTA, message, NOCOLOR)
}

func underlineOutput(message string) string {
	return fmt.Sprintf("%v%v%v", UNDERLINE, message, NOCOLOR)
}
