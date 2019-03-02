package main

import (
	"fmt"
	"github.com/codesoap/ytools/src/ytools"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const max_results = 12

type Video struct {
	Title string
	Url   string
}

func main() {
	search_url := get_search_url()
	videos, err := scrape_off_videos(search_url)
	if err != nil {
		os.Exit(1)
	}
	if len(videos) == 0 {
		fmt.Fprintf(os.Stderr, "No videos found.\n")
		os.Exit(1)
	}
	if err := save_urls(videos); err != nil {
		fmt.Fprintf(os.Stderr, "Failed saving found URLs.\n")
		os.Exit(1)
	}
	print_video_titles(videos)
}

func get_search_url() string {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Give one or more search terms as parameters.\n")
		os.Exit(1)
	}
	search_string := url.PathEscape(strings.Join(os.Args[1:], " "))
	return fmt.Sprintf(
		"https://www.youtube.com/results?search_query=%s",
		search_string)
}

func scrape_off_videos(search_url string) (videos []Video, err error) {
	videos = make([]Video, 0, max_results)

	resp, err := http.Get(search_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get '%s'\n", search_url)
		return
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenizer_status := tokenizer.Next()
		if tokenizer_status == html.ErrorToken {
			break
		} else if tokenizer_status == html.StartTagToken {
			if is_result(tokenizer.Token()) {
				video, ok := extract_video(tokenizer)
				if !ok {
					break
				}
				videos = append(videos, video)
				if len(videos) == max_results {
					break
				}
			}
		}
	}
	return
}

func save_urls(videos []Video) (err error) {
	data_dir, err := ytools.GetDataDir()
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
		err = urls_file.Close()
	}()
	for _, video := range videos {
		_, err := fmt.Fprintf(urls_file, "https://www.youtube.com%s\n", video.Url)
		if err != nil {
			return err
		}
	}
	return
}

func print_video_titles(videos []Video) {
	for i, video := range videos {
		fmt.Printf("%2d: %s\n", i+1, video.Title)
	}
}

func is_result(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "class" && strings.Contains(a.Val, "yt-lockup-video") {
			return true
		}
	}
	return false
}

func extract_video(tokenizer *html.Tokenizer) (video Video, ok bool) {
	for {
		tokenizer_status := tokenizer.Next()
		if tokenizer_status == html.ErrorToken {
			return
		} else if tokenizer_status == html.StartTagToken {
			token := tokenizer.Token()
			if is_title_link(token) {
				video, ok = extract_video_from_title_link(token)
				return
			}
		}
	}
}

func is_title_link(token html.Token) bool {
	if token.Data == "a" {
		for _, a := range token.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "yt-uix-tile-link") {
				return true
			}
		}
	}
	return false
}

func extract_video_from_title_link(token html.Token) (video Video, ok bool) {
	for _, a := range token.Attr {
		if a.Key == "title" {
			video.Title = a.Val
		}
		if a.Key == "href" {
			video.Url = a.Val
		}
	}
	if video.Title != "" && video.Url != "" {
		ok = true
	}
	return
}
