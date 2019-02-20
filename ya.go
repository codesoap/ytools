package main

import (
	"os"
	"os/exec"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const result_cnt = 16
var api_url = "http://youtube-scrape.herokuapp.com/api/search?q=%s&page=1"

type Scrape struct {
	Results []struct {
		Video struct {
			Title string
			Url   string
		}
	}
}

func main() {
	// TODO: If stdin is given: use it as search terms
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Give one or more search terms as parameters.")
		os.Exit(1)
	}
	scrape := get_search_results(strings.Join(os.Args[1:], "%20"))
	selection := get_selection_from_user(scrape)
	url := scrape.Results[selection - 1].Video.Url
	play_audio(url)
}

func get_search_results(search string) Scrape {
	url := fmt.Sprintf(api_url, search)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var scrape Scrape
	json.Unmarshal(body, &scrape)

	return scrape
}

func get_selection_from_user(scrape Scrape) int {
    for i, result := range scrape.Results {
		if i >= result_cnt {
			break
		}
		fmt.Printf("%2d: %s\n", i + 1, result.Video.Title)
	}
	var selection int
	fmt.Print("Selection: ")
	if _, err := fmt.Scanf("%d", &selection); err != nil {
		panic(err)
	}
	if selection < 1 || selection > result_cnt || selection > len(scrape.Results) {
		fmt.Fprintln(os.Stderr, "Selection out of range.")
		os.Exit(1)
	}
	return selection
}

func play_audio(url string) {
	cmd := exec.Command("mpv",
	                    "--ytdl-format", "bestaudio/best",
	                    "--no-video",
	                    url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
	    panic(err)
	}
	defer func() {
		if err := cmd.Wait(); err != nil {
		    panic(err)
		}
	}()
}
