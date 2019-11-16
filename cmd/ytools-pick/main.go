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
		os.Exit(1)
	}
	save_as_last_picked(url)
	fmt.Println(url)
}

func save_as_last_picked(url string) (err error) {
	data_dir, err := ytools.GetDataDir()
	if err != nil {
		return
	}
	last_picked_filename := filepath.Join(data_dir, "last_picked")
	last_picked_file, err := os.Create(last_picked_filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create last_picked file.")
		return
	}
	defer func() {
		err = last_picked_file.Close()
	}()
	_, err = fmt.Fprintln(last_picked_file, url)
	if err != nil {
		return
	}
	return
}
