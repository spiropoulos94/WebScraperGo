package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	// url := "http://example.com" // replace with your URL
	url := "https://preview.themeforest.net/item/edumate-education-html-template/full_screen_preview/27878270?_ga=2.195597292.139702443.1673861451-242333065.1673861451" // replace with your URL

	// Get the HTML from the URL
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Create a file to save the HTML
	htmlFile, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer htmlFile.Close()

	// Copy the HTML from the response to the file
	_, err = io.Copy(htmlFile, resp.Body)
	if err != nil {
		panic(err)
	}

	// Read the HTML file into a string
	htmlBytes, err := os.Open("index.html")
	if err != nil {
		panic(err)
	}
	defer htmlBytes.Close()

	html := make([]byte, resp.ContentLength)
	htmlBytes.Read(html)

	// Replace the URLs of linked assets with local file paths
	r, _ := regexp.Compile("(?i)https?://[^/\\s]+")

	css := r.FindAllString(string(html), -1)
	js := r.FindAllString(string(html), -1)
	img := r.FindAllString(string(html), -1)

	for _, link := range css {
		if !strings.HasSuffix(link, ".css") {
			continue
		}

		// Download the CSS file
		resp, err := http.Get(link)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		// Create a file to save the CSS
		fileName := strings.TrimPrefix(link, "http://")
		fileName = strings.TrimPrefix(fileName, "https://")
		fileName = strings.Replace(fileName, "/", "_", -1)
		cssFile, err := os.Create(fileName)
		if err != nil {
			continue
		}
		defer cssFile.Close()

		// Copy the CSS from the response to the file
		_, err = io.Copy(cssFile, resp.Body)
		if err != nil {
			continue
		}

		// Replace the URL in the HTML with the local file path
		html = []byte(strings.Replace(string(html), link, fileName, -1))
	}
	for _, link := range js {
		if !strings.HasSuffix(link, ".js") {
			continue
		}

		// Download the JS file
		resp, err := http.Get(link)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		// Create a file to save the JS
		fileName := strings.TrimPrefix(link, "http://")
		fileName = strings.TrimPrefix(fileName, "https://")
		fileName = strings.Replace(fileName, "/", "_", -1)
		jsFile, err := os.Create(fileName)
		if err != nil {
			continue
		}
		defer jsFile.Close()

		// Copy the JS from the response to the file
		_, err = io.Copy(jsFile, resp.Body)
		if err != nil {
			continue
		}

		// Replace the URL in the HTML with the local file path
		html = []byte(strings.Replace(string(html), link, fileName, -1))
	}
	for _, link := range img {
		if !strings.HasSuffix(link, ".jpg") && !strings.HasSuffix(link, ".jpeg") && !strings.HasSuffix(link, ".png") && !strings.HasSuffix(link, ".gif") && !strings.HasSuffix(link, ".svg") {
			continue
		}

		// Download the image file
		resp, err := http.Get(link)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		// Create a file to save the image
		fileName := strings.TrimPrefix(link, "http://")
		fileName = strings.TrimPrefix(fileName, "https://")
		fileName = strings.Replace(fileName, "/", "_", -1)
		imgFile, err := os.Create(fileName)
		if err != nil {
			continue
		}
		defer imgFile.Close()

		// Copy the image from the response to the file
		_, err = io.Copy(imgFile, resp.Body)
		if err != nil {
			continue
		}

		// Replace the URL in the HTML with the local file path
		html = []byte(strings.Replace(string(html), link, fileName, -1))
	}

	// Save the updated HTML with local file paths
	htmlFile, err = os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer htmlFile.Close()

	htmlFile.Write(html)
}

