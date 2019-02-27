package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr,
			"Please give the search result number as an argument.")
		os.Exit(1)
	}
	search_results, err := get_search_results()
	if err != nil {
		os.Exit(1)
	}
	selection, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "The given argument is no integer.")
		os.Exit(1)
	}
	if selection < 1 || selection > len(search_results) {
		fmt.Fprintln(os.Stderr, "Selection out of range.")
		os.Exit(1)
	}
	url := search_results[selection-1]
	save_as_last_picked(url)
	fmt.Println(url)
}

func get_search_results() (search_results []string, err error) {
	search_results = make([]string, 0)

	data_dir := get_data_dir()
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

func save_as_last_picked(url string) (err error) {
	data_dir := get_data_dir()
	err = os.MkdirAll(data_dir, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create directory '%s'.", data_dir)
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

func get_data_dir() string {
	data_dir := os.Getenv("XDG_DATA_HOME")
	if data_dir == "" {
		data_dir = filepath.Join(os.Getenv("HOME"), ".local/share/ytools/")
	}
	return data_dir
}
