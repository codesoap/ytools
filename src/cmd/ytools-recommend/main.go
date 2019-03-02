package main

import (
	"fmt"
	"github.com/codesoap/ytools/src/ytools"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

const max_results = 16

type Video struct {
	Title string
	Url   string
}

func main() {
	video_url, err := ytools.GetDesiredVideoUrl()
	if err != nil {
		os.Exit(1)
	}
	recommendations, err := scrape_off_recommendations(video_url)
	if err != nil {
		os.Exit(1)
	}
	if len(recommendations) == 0 {
		fmt.Fprintf(os.Stderr, "No recommendations found.\n")
		os.Exit(1)
	}
	if err := save_recommendations_urls(recommendations); err != nil {
		fmt.Fprintf(os.Stderr, "Failed saving found URLs.\n")
		os.Exit(1)
	}
	print_video_titles(recommendations)
}

func scrape_off_recommendations(video_url string) (videos []Video, err error) {
	videos = make([]Video, 0, max_results)

	resp, err := http.Get(video_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get '%s'\n", video_url)
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
			if is_recommendation(token) {
				video, ok := extract_next_video(token)
				if !ok {
					err = fmt.Errorf("failed at video extraction")
					return
				}
				videos = append(videos, video)
				if len(videos) >= max_results {
					return
				}
			}
		}
	}
	return
}

func save_recommendations_urls(videos []Video) (err error) {
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

func is_recommendation(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "class" && strings.Contains(a.Val, "content-link") {
			return true
		}
	}
	return false
}

func extract_next_video(token html.Token) (video Video, ok bool) {
	for _, a := range token.Attr {
		if a.Key == "href" {
			video.Url = fmt.Sprintf("https://www.youtube.com%s", a.Val)
		} else if a.Key == "title" {
			video.Title = a.Val
		}
	}
	if video.Url != "" && video.Title != "" {
		ok = true
	}
	return
}
