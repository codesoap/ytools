package main

import (
	"fmt"
	"github.com/codesoap/ytools"
	"os"
	"path/filepath"
)

func main() {
	url, err := ytools.GetDesiredVideoUrl()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get video URL: %s\n", err.Error())
		os.Exit(1)
	}
	saveAsLastPicked(url)
	fmt.Println(url)
}

func saveAsLastPicked(url string) (err error) {
	dataDir, err := ytools.GetDataDir()
	if err != nil {
		return
	}
	lastPickedFilename := filepath.Join(dataDir, "last_picked")
	lastPickedFile, err := os.Create(lastPickedFilename)
	if err != nil {
		return
	}
	defer func() {
		err = lastPickedFile.Close()
	}()
	_, err = fmt.Fprintln(lastPickedFile, url)
	if err != nil {
		return
	}
	return
}
