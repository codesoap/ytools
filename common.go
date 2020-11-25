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
	dataDir, err := GetDataDir()
	if err != nil {
		return
	}
	urlsFilename := filepath.Join(dataDir, "search_results")
	urlsFile, err := os.Create(urlsFilename)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			urlsFile.Close()
		} else {
			err = urlsFile.Close()
		}
	}()
	for _, url := range urls {
		_, err = fmt.Fprintln(urlsFile, url)
	}
	return
}

func GetDesiredVideoUrl() (videoUrl string, err error) {
	switch len(os.Args) {
	case 1:
		videoUrl, err = GetLastPickedUrl()
	case 2:
		var selection int
		selection, err = strconv.Atoi(os.Args[1])
		if err != nil {
			return
		}
		videoUrl, err = GetSearchResult(selection - 1)
	default:
		err = fmt.Errorf("invalid argument count; give a video number as " +
			"argument, or no argument to select the last picked")
	}
	return
}

func GetSearchResult(i int) (searchResult string, err error) {
	searchResults, err := getSearchResults()
	if err == nil {
		if i < 0 || i >= len(searchResults) {
			err = fmt.Errorf("invalid search result index")
		} else {
			searchResult = searchResults[i]
		}
	}
	return
}

func getSearchResults() (searchResults []string, err error) {
	searchResults = make([]string, 0)

	dataDir, err := GetDataDir()
	if err != nil {
		return
	}
	urlsFile := filepath.Join(dataDir, "search_results")
	file, err := os.Open(urlsFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		searchResults = append(searchResults, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return
}

func GetLastPickedUrl() (lastPickedUrl string, err error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return
	}
	lastPickedFilename := filepath.Join(dataDir, "last_picked")
	fileContent, err := ioutil.ReadFile(lastPickedFilename)
	if err != nil {
		return
	}
	lastPickedUrl = strings.TrimSpace(string(fileContent))
	return
}

func GetDataDir() (dataDir string, err error) {
	dataDirBase := os.Getenv("XDG_DATA_HOME")
	if dataDirBase == "" {
		dataDirBase = filepath.Join(os.Getenv("HOME"), ".local/share/")
	}
	dataDir = filepath.Join(dataDirBase, "ytools/")
	err = os.MkdirAll(dataDir, 0755)
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
	re := regexp.MustCompile(`(?m)ytInitialData.{1,3}= *(.*?);(</script>|$)`)
	matches := re.FindSubmatch(bytes)
	if matches == nil {
		err = fmt.Errorf("retrieved HTML does not contain the expected JSON")
		return
	}
	return matches[1], nil
}
