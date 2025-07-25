package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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

	if _, err := file.Write(audioData); err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(file.Name()); err != nil {
			log.Printf("failed to remove file %s: %v", file.Name(), err)
		}
	}()
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s: %v", file.Name(), err)
		}
	}()

	switch runtime.GOOS {
	case "darwin":
		if err := exec.Command("afplay", file.Name()).Run(); err != nil {
			return err
		}
	case "linux":
		if err := exec.Command("aplay", file.Name()).Run(); err != nil {
			return err
		}
	case "windows":
		cmd := fmt.Sprintf("powershell -c (New-Object Media.SoundPlayer '%s').PlaySync();", file.Name())
		if err := exec.Command(cmd).Run(); err != nil {
			return err
		}
	}
	return nil
}
