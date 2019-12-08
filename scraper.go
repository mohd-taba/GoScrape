package scraper

/*
Open a series of urls.
Check status code for each url and store urls I could not
open in a dedicated array.
Fetch urls concurrently using goroutines.
*/

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// Options ... Configuration of the created scrape
type Options struct {
	// URLSlice ... (Required) A list of strings of the urls of type []string
	URLSlice []string
	// ProxyURL ... (Optional) (e.g. "http://localhost:8118")
	ProxyURL string
	// UserAgent ... (Optional) Set custom user-agent
	UserAgent string
	// Jar ... (Optional) Pass your own cookiejar after signing in for example
	Jar *cookiejar.Jar
	//Callback ... (Optional) I know I just met you, but call me maybe
	CallbackF func(resp *http.Response)
}
// Default Values.
var (
	userAgent = "GoScrape"
)
// fetchURL opens a url with GET method and sets a custom user agent.
// If url cannot be opened, then log it to a dedicated channel.
func fetchURL(uri string, proxy string, cookieJar *cookiejar.Jar, chFailedUrls chan string, chIsFinished chan *http.Response) {
	//Preparing proxy
	var client = &http.Client{}
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		client = &http.Client{Transport: transport, Jar: cookieJar}
		if err != nil {
			panic("Proxy Error")
		}
	}

	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)

	// Inform the channel chIsFinished that url fetching is done (no
	// matter whether successful or not). Defer triggers only once
	// we leave fetchURL():
	defer func() {
		chIsFinished <- resp
	}()

	// If url could not be opened, we inform the channel chFailedUrls:
	if err != nil || resp.StatusCode != 200 {
		chFailedUrls <- uri
		resp.Body.Close()
		return
	}

}

// Scrape ... Takes scraper.Options as argument, and scrapes provided urls accordingly
func Scrape(opts Options) {
	// Process options
	if len(opts.URLSlice) == 0{
		panic("Empty URL list passed")
	}
	if opts.UserAgent != "" {
		userAgent = opts.UserAgent
	}
	if opts.CallbackF == nil{
		opts.CallbackF = func(resp *http.Response) {

		}

	}
	// Create 2 channels, 1 to track urls we could not open
	// and 1 to inform url fetching is done:
	chFailedUrls := make(chan string)
	chIsFinished := make(chan *http.Response)


	// Open all urls concurrently using the 'go' keyword:
	for _, uri := range opts.URLSlice {
		go fetchURL(uri, opts.ProxyURL, opts.Jar, chFailedUrls, chIsFinished)
	}

	// Receive messages from every concurrent goroutine. If
	// an url fails, we log it to failedUrls array:
	failedUrls := make([]string, 0)
	for i := 0; i < len(opts.URLSlice); {
		select {
		case uri := <-chFailedUrls:
			failedUrls = append(failedUrls, uri)
		case resp := <-chIsFinished:
			opts.CallbackF(resp)
			defer resp.Body.Close()
			i++
		}
	}

	// Print all urls we could not open:
	fmt.Println("Could not fetch these urls: ", failedUrls)

}
