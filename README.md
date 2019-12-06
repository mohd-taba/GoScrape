GoScrape is a simple web scraper written in Go.
It includes proxy support, user-agent, and cookies.

# Installation
`go get github.com/mohd-taba/GoScrape`

# Usage

```golang
// Define your Options
	opts := scraper.Options{
		URLSlice : []string {"https://duckduckgo.com/", "https://example.com/", "http://site1.com/"},
		UserAgent : "GoScrape",
    		ProxyURL : "http://127.0.0.1:8118"
		Jar : MyCookieJar,
	}
  
  //Run Scraper
  scraper.Scrape(opts)
  
  //That's it!
  ```
  
  ## Cookiejar?
  
  ```golang
  // Create cookie jar:
  options := cookiejar.Options{
        PublicSuffixList: publicsuffix.List,
    }
    jar, err := cookiejar.New(&options)
    if err != nil {
        log.Fatal(err)
    }
    // Do stuff with cookie jar
    client := http.Client{Jar: jar}
    resp, err := client.Get("http://dubbelboer.com/302cookie.php")
    if err != nil {
        log.Fatal(err)
    }
    //Create options with jar
    opts = scraper.Options{
    URLSlice : []string {"https://cookiesite.com"},
    }
    scraper.Scrape(opts)
    ```
  
