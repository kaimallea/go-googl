package main

import (
	"http"
	"fmt"
	"io/ioutil"
	"json"
	"os"
	"strings"
)


var (
	expand_endpoint  string = "https://www.googleapis.com/urlshortener/v1/url?shortUrl="
	shorten_endpoint string = "https://www.googleapis.com/urlshortener/v1/url"
)


/**
 * Expand a shortened url
 *
 */
func expandUrl(url string) (string, string) {
	res, _, err := http.Get(expand_endpoint + url)
	if err != nil {
		fmt.Printf("%s", err)
		return "", url
	}

	var jsonresult map[string]string

	body, _ := ioutil.ReadAll(res.Body)
	if json.Unmarshal(body, &jsonresult) != nil {
		fmt.Printf("Error processing %s", url)
		return "", url
	}

	res.Body.Close()

	return jsonresult["longUrl"], url
}


/**
 * Shorten a long url
 *
 */
func shortenUrl(url string) (string, string) {
	payload := strings.NewReader("{\"longUrl\": \"" + url + "\"}")

	res, err := http.Post(shorten_endpoint, "application/json", payload)
	if err != nil {
		fmt.Printf("%s", err)
		return "", url
	}

	var jsonresult map[string]string

	body, _ := ioutil.ReadAll(res.Body)
	if json.Unmarshal(body, &jsonresult) != nil {
		fmt.Printf("Error processing %s", url)
		return "", url
	}

	res.Body.Close()

	return jsonresult["id"], url
}


/**
 * Process a url -- determine if it should be
 * shortened or expanded
 *
 */
func processUrl(url string) (string, string) {

	// Prepend protocol if its missing
	if !strings.Contains(url, "http://") {
		url = "http://" + url
	}

	// Expand url if it's already shortened
	if strings.Contains(url, "goo.gl") {
		return expandUrl(url)
	}

	return shortenUrl(url)
}


func main() {
	nArgs := len(os.Args)

	if nArgs < 2 {
		fmt.Println("Specify some URLs")
		return
	}

	for i := 1; i < nArgs; i++ {
		r, u := processUrl(os.Args[i])
		fmt.Println(u, " -> ", r)
	}
}
