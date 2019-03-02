package main

import (
	"fmt"
	"github.com/codesoap/ytools/src/ytools"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

type Info struct {
	Title string
	Views string
	// TODO: Length string
	// The length is not available on the main page, will probably have
	// to load something like this:
	// https://www.youtube.com/annotations_invideo?video_id=DuoTdnq_OqE
	Likes       string
	Dislikes    string
	Date        string
	Description string
}

func main() {
	video_url, err := ytools.GetDesiredVideoUrl()
	if err != nil {
		os.Exit(1)
	}
	info, err := scrape_off_info(video_url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to scrape the videos page")
		os.Exit(1)
	}
	print_info(info)
}

func print_info(info Info) {
	fmt.Println(info.Title)
	fmt.Printf("%s  ▲ %s  ▼ %s  %s\n\n", info.Views, info.Likes, info.Dislikes,
		info.Date)
	fmt.Println(info.Description)
}

func scrape_off_info(url string) (info Info, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			err = fill_in_available_info(tokenizer, &info)
			if err != nil {
				return
			}
		}
	}
	if !info_complete(info) {
		err = fmt.Errorf("missing some info field(s)")
	}
	return
}

func fill_in_available_info(tokenizer *html.Tokenizer, info *Info) (err error) {
	ok := true
	token := tokenizer.Token()
	switch {
	case is_title(token):
		info.Title, ok = extract_next_text(tokenizer)
	case is_views(token):
		info.Views, ok = extract_next_text(tokenizer)
	case is_likes(token):
		info.Likes, ok = extract_next_text(tokenizer)
	case is_dislikes(token):
		info.Dislikes, ok = extract_next_text(tokenizer)
	case is_date(token):
		info.Date, ok = extract_next_text(tokenizer)
	case is_description(token):
		info.Description, ok = extract_description(tokenizer)
	}
	if !ok {
		err = fmt.Errorf("failed at extraction")
	}
	return
}

func is_title(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "class" && a.Val == "watch-title" {
			return true
		}
	}
	return false
}

func is_views(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "class" && a.Val == "watch-view-count" {
			return true
		}
	}
	return false
}

// FIXME: Results are off by one
func is_likes(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "class" &&
			strings.Contains(a.Val, "like-button-renderer-like-button") {
			return true
		}
	}
	return false
}

// FIXME: Results are off by one
func is_dislikes(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "class" &&
			strings.Contains(a.Val, "like-button-renderer-dislike-button") {
			return true
		}
	}
	return false
}

func is_date(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "class" && a.Val == "watch-time-text" {
			return true
		}
	}
	return false
}

func is_description(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "id" && a.Val == "eow-description" {
			return true
		}
	}
	return false
}

func extract_description(tokenizer *html.Tokenizer) (desc string, ok bool) {
	description := make([]byte, 100)
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return
		case html.TextToken:
			description = append(description, tokenizer.Text()...)
		case html.SelfClosingTagToken:
			if tokenizer.Token().Data == "br" {
				description = append(description, '\n')
			}
		case html.StartTagToken:
			if tokenizer.Token().Data == "a" {
				var next_text string
				next_text, ok = extract_next_text(tokenizer)
				if !ok {
					return
				}
				description = append(description, []byte(next_text)...)
			}
		case html.EndTagToken:
			if tokenizer.Token().Data == "p" {
				return string(description), true
			}
		}
	}
}

func extract_next_text(tokenizer *html.Tokenizer) (text string, ok bool) {
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return
		case html.TextToken:
			raw_text := tokenizer.Text()
			tmp_text := make([]byte, len(raw_text))
			copy(tmp_text, raw_text)
			return strings.TrimSpace(string(tmp_text)), true
		}
	}
}

func info_complete(info Info) bool {
	return info.Title != "" &&
		info.Views != "" &&
		info.Likes != "" &&
		info.Dislikes != "" &&
		info.Date != "" &&
		info.Description != ""
}
