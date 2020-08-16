package ytools

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
		return
	}
	defer func() {
		if err != nil {
			urls_file.Close()
		} else {
			err = urls_file.Close()
		}
	}()
	for _, url := range urls {
		_, err = fmt.Fprintln(urls_file, url)
	}
	return
}

func GetDesiredVideoUrl() (video_url string, err error) {
	switch len(os.Args) {
	case 1:
		video_url, err = GetLastPickedUrl()
	case 2:
		var selection int
		selection, err = strconv.Atoi(os.Args[1])
		if err != nil {
			return
		}
		video_url, err = GetSearchResult(selection - 1)
	default:
		err = fmt.Errorf("invalid argument count; give a video number as " +
			"argument, or no argument to select the last picked")
	}
	return
}

func GetSearchResult(i int) (search_result string, err error) {
	search_results, err := get_search_results()
	if err == nil {
		if i < 0 || i >= len(search_results) {
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
		return
	}
	last_picked_url = strings.TrimSpace(string(file_content))
	return
}

func GetDataDir() (data_dir string, err error) {
	data_dir_base := os.Getenv("XDG_DATA_HOME")
	if data_dir_base == "" {
		data_dir_base = filepath.Join(os.Getenv("HOME"), ".local/share/")
	}
	data_dir = filepath.Join(data_dir_base, "ytools/")
	err = os.MkdirAll(data_dir, 0755)
	return
}

// ExtractJson returns the ytInitialData JSON from the HTML at the
// given URL.
func ExtractJson(url string) (mainJson []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	re := regexp.MustCompile(`(?m)^ *window\["ytInitialData"\] *= *(.*); *$`)
	matches := re.FindSubmatch(bytes)
	if matches == nil {
		err = fmt.Errorf("retrieved HTML does not contain the expected JSON")
		return
	}
	return matches[1], nil
}
