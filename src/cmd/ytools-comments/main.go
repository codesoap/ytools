package main

import (
	"encoding/json"
	"fmt"
	"github.com/codesoap/ytools/src/ytools"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const max_results = 4

var comments_url_base string = "https://www.youtube.com/" +
	"comment_service_ajax?action_get_comments=1&ctoken=%s"

type CommentsJson struct {
	Content_html string
}

func main() {
	video_url, err := ytools.GetDesiredVideoUrl()
	if err != nil {
		os.Exit(1)
	}
	comments_html, err := get_comments_html(video_url)
	if err != nil {
		os.Exit(1)
	}
	comments, ok := extract_comments(comments_html)
	if !ok {
		fmt.Fprintf(os.Stderr, "Extracting the comments failed\n")
		os.Exit(1)
	}
	print_comments(comments)
}

func get_comments_html(video_url string) (comments_html string, err error) {
	comments_req, err := get_comments_request(video_url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(comments_req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get the comments\n")
		return
	}
	json_bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read body of the comments response\n")
		return
	}
	var comments_json CommentsJson
	err = json.Unmarshal(json_bytes, &comments_json)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not unmarshal the comments JSON\n")
	}
	return comments_json.Content_html, err
}

func extract_comments(comments_html string) (comments []string, ok bool) {
	var comment string
	ok = true
	comments = make([]string, 0, max_results)
	tokenizer := html.NewTokenizer(strings.NewReader(comments_html))
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return
			} else {
				return nil, false
			}
		case html.StartTagToken:
			if is_comment(tokenizer.Token()) {
				comment, ok = extract_comment(tokenizer)
				if !ok {
					return nil, false
				}
				comments = append(comments, comment)
				if len(comments) >= max_results {
					return
				}
			}
		}
	}
	return
}

func print_comments(comments []string) {
	for i, comment := range comments {
		fmt.Printf("=== Comment #%d: ===\n", i+1)
		fmt.Println(comment, "\n")
	}
}

func is_comment(token html.Token) bool {
	for _, a := range token.Attr {
		if a.Key == "class" && a.Val == "comment-renderer-text-content" {
			return true
		}
	}
	return false
}

func extract_comment(tokenizer *html.Tokenizer) (comment string, ok bool) {
	comment_bytes := make([]byte, 0, 100)
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return
		case html.SelfClosingTagToken:
			if tokenizer.Token().Data == "br" {
				comment_bytes = append(comment_bytes, '\n')
			}
		case html.TextToken:
			comment_bytes = append(comment_bytes, tokenizer.Text()...)
		case html.EndTagToken:
			if tokenizer.Token().Data == "div" {
				return string(comment_bytes), true
			}
		}
	}
}

// Compiling the request for the comments with cookies and session_token
// is cumbersome, but required to avoid a 403 error code.
func get_comments_request(video_url string) (req *http.Request, err error) {
	resp, err := http.Get(video_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get '%s'\n", video_url)
		return
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read body of '%s'\n", video_url)
		return
	}
	video_html := string(buf)

	req_url, ok := extract_request_url(video_html)
	if !ok {
		return nil, fmt.Errorf("Could not extract comments URL")
	}
	req_body_reader, ok := get_request_body_reader(video_html)
	if !ok {
		return nil, fmt.Errorf("Could not extract session token")
	}
	req, err = http.NewRequest("POST", req_url, req_body_reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create request for the comments\n")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}

	return
}

func extract_request_url(text string) (url string, ok bool) {
	ctoken, ok := extract_string_after_key(text, "COMMENTS_TOKEN")
	return fmt.Sprintf(comments_url_base, ctoken), ok
}

func get_request_body_reader(text string) (reader io.Reader, ok bool) {
	xsrf_token, ok := extract_string_after_key(text, "XSRF_TOKEN")
	body := fmt.Sprintf("session_token=%s", xsrf_token)
	return strings.NewReader(body), ok
}

func extract_string_after_key(text string, key string) (s string, ok bool) {
	i_key := strings.Index(text, key)
	if i_key < 0 {
		return
	}
	i_start := strings.IndexByte(text[i_key:], '"') + 1 + i_key
	i_end := strings.IndexByte(text[i_start:], '"') + i_start
	return text[i_start:i_end], true
}
