package ytools

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func GetSearchResults() (search_results []string, err error) {
	search_results = make([]string, 0)

	data_dir, err := GetDataDir()
	if err != nil {
		return
	}
	urls_file := filepath.Join(data_dir, "search_results")
	file, err := os.Open(urls_file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read '%s'.", urls_file)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		search_results = append(search_results, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return
}

func GetDataDir() (data_dir string, err error) {
	data_dir = os.Getenv("XDG_DATA_HOME")
	if data_dir == "" {
		data_dir = filepath.Join(os.Getenv("HOME"), ".local/share/ytools/")
	}
	err = os.MkdirAll(data_dir, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create directory '%s'.", data_dir)
	}
	return
}
