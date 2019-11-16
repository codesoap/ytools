package main

import (
	"fmt"
	"github.com/codesoap/ytools"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
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
	if err := save_videos_urls(videos); err != nil {
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
		switch tokenizer.Next() {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			token := tokenizer.Token()
			if is_video_title_link(token) {
				video, ok := extract_video_from_title_link(token)
				if !ok {
					return
				}
				videos = append(videos, video)
				if len(videos) == max_results {
					return
				}
			}
		}
	}
}

func save_videos_urls(videos []Video) (err error) {
	videos_urls := make([]string, 0, max_results)
	for _, video := range videos {
		videos_urls = append(videos_urls, video.Url)
	}
	return ytools.SaveUrls(videos_urls)
}

func print_video_titles(videos []Video) {
	for i, video := range videos {
		fmt.Printf("%2d: %s\n", i+1, video.Title)
	}
}

func is_video_title_link(token html.Token) bool {
	is_tile_link, is_video := false, false
	if token.Data == "a" {
		for _, a := range token.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "yt-uix-tile-link") {
				is_tile_link = true
			}
			// This filters out channels and playlists:
			if a.Key == "href" && strings.HasPrefix(a.Val, "/watch") &&
				len(a.Val) == 20 {
				is_video = true
			}
		}
	}
	return is_tile_link && is_video
}

func extract_video_from_title_link(token html.Token) (video Video, ok bool) {
	for _, a := range token.Attr {
		if a.Key == "title" {
			video.Title = a.Val
		}
		if a.Key == "href" {
			video.Url = fmt.Sprintf("https://www.youtube.com%s", a.Val)
		}
	}
	if video.Title != "" && video.Url != "" {
		ok = true
	}
	return
}
