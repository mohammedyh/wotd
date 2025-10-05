package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func playPronunciationAudio() error {
	const audioURLPrefix = "https://media.merriam-webster.com/audio/prons/en/us/mp3"
	doc, err := initDocumentReader(wotdURL)
	if err != nil {
		return fmt.Errorf("error parsing html document: %v", err)
	}

	filename, exists := doc.Find(".word-and-pronunciation .play-pron").Attr("data-file")
	if !exists {
		return errors.New("unable to get audio filename")
	}
	audioURL := fmt.Sprintf("%v/%v/%v.mp3", audioURLPrefix, string(filename[0]), filename)

	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Get(audioURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error %d %s", res.StatusCode, res.Status)
	}

	audioData, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read audio file")
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get user home directory")
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
