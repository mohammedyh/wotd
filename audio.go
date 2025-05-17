package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func playProunciationAudio() error {
	doc, err := initDocumentReader(WOTD_URL)
	if err != nil {
		return fmt.Errorf("error parsing html document: %v", err)
	}

	audioFileUrl, audioFileExists := doc.Find(".otd-item-headword__pronunciation-audio").Attr("href")
	if !audioFileExists {
		return fmt.Errorf("no pronunciation audio for this word\n")
	}

	res, err := http.Get(audioFileUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error %d %s\n", res.StatusCode, res.Status)
	}

	audioData, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read audio file")
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("can't get user home directory")
	}

	file, err := os.CreateTemp(homedir, "pronunciation-audio*.mp3")
	if err != nil {
		return err
	}

	_, err = file.Write(audioData)
	if err != nil {
		return err
	}

	defer os.Remove(file.Name())
	defer file.Close()

	cmd := exec.Command("afplay", file.Name())
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
