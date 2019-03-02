package ytools

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func SaveUrls(urls []string) (err error) {
	data_dir, err := GetDataDir()
	if err != nil {
		return
	}
	urls_filename := filepath.Join(data_dir, "search_results")
	urls_file, err := os.Create(urls_filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create URLs file.")
		return
	}
	defer func() {
		// FIXME: This overwrites previous errors
		err = urls_file.Close()
	}()
	for _, url := range urls {
		_, err = fmt.Fprintln(urls_file, url)
		if err != nil {
			return
		}
	}
	return
}

func GetDesiredVideoUrl() (video_url string, err error) {
	switch len(os.Args) {
	case 1:
		video_url, err = GetLastPickedUrl()
	case 2:
		selection, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "The given argument is no integer.")
		}
		video_url, err = GetSearchResult(selection - 1)
	default:
		fmt.Fprintf(os.Stderr, "Give a video number as argument, or no "+
			"argument to select the last picked.\n")
		err = fmt.Errorf("invalid argument count")
	}
	return
}

func GetSearchResult(i int) (search_result string, err error) {
	search_results, err := get_search_results()
	if err == nil {
		if i < 0 || i >= len(search_results) {
			fmt.Fprintln(os.Stderr, "Search result index out of range.")
			err = fmt.Errorf("invalid search result index")
		} else {
			search_result = search_results[i]
		}
	}
	return
}

func get_search_results() (search_results []string, err error) {
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

func GetLastPickedUrl() (last_picked_url string, err error) {
	data_dir, err := GetDataDir()
	if err != nil {
		return
	}
	last_picked_filename := filepath.Join(data_dir, "last_picked")
	file_content, err := ioutil.ReadFile(last_picked_filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read '%s'.", last_picked_filename)
		return
	}
	last_picked_url = strings.TrimSpace(string(file_content))
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
