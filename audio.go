package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func playProunciationAudio() (string, error) {
	doc, err := initDocumentReader(WOTD_URL)
	if err != nil {
		return "", fmt.Errorf("error parsing html document: %v", err)
	}

	audioFileUrl, audioFileExists := doc.Find(".otd-item-headword__pronunciation-audio").Attr("href")
	if !audioFileExists {
		return "", fmt.Errorf("no pronounciation audio for this word\n")
	}

	res, err := http.Get(audioFileUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error %d %s\n", res.StatusCode, res.Status)
	}

	data, _ := io.ReadAll(res.Body)
	home, _ := os.UserHomeDir()
	file, err := os.CreateTemp(home, "pronounciation-audio*.mp3")
	if err != nil {
		return "", err
	}

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	defer file.Close()
	defer os.Remove(file.Name())

	cmd := exec.Command("afplay", file.Name())
	cmd.Run()
	return string(data), nil
}
