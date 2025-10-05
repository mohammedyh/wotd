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
	wotdURL   = "https://www.merriam-webster.com/word-of-the-day/"
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
		open(wotdURL)
		os.Exit(0)
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
		word              = doc.Find(".word-header-txt")
		wordType          = doc.Find(".word-attributes .main-attr")
		definition        = doc.Find(".wod-definition-container p").First()
		examples          = doc.Find(".wod-definition-container p").Slice(1, -5).Text()
		formattedExamples = strings.TrimSpace(strings.Join(strings.Split(examples, "//"), "\n-"))
	)

	fmt.Printf("%v: %v\n", colorOutput("Word"), word.First().Text())
	fmt.Printf("%v: %v\n", underlineOutput("Word Type"), strings.Trim(wordType.Text(), " \n"))
	fmt.Printf("%v: %v\n", underlineOutput("Definition"), definition.Text())
	fmt.Printf("%v:\n%v\n", underlineOutput("Examples"), formattedExamples)
}

func open(url string) {
	var program string

	switch runtime.GOOS {
	case "linux":
		program = "xdg-open"
	case "darwin":
		program = "open"
	case "windows":
		program = "start"
	}
	if err := exec.Command(program, url).Run(); err != nil {
		log.Fatal("couldn't open word of the day page in browser")
	}
}

func colorOutput(message string) string {
	return fmt.Sprintf("%v %v %v", magenta, message, noColor)
}

func underlineOutput(message string) string {
	return fmt.Sprintf("%v%v%v", underline, message, noColor)
}
