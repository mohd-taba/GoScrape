# Description
GoScrape is a simple web scraper written in Go.
It supports concurrency, proxy, user-agent, and cookies.

# Installation
`go get github.com/mohd-taba/GoScrape`

# Usage

WARNING: In case you pass a callback function, don't forget to resp.Body.Close()

```golang
// Define your Options
	opts := scraper.Options{
		URLSlice : []string {"https://duckduckgo.com/", "https://example.com/", "http://site1.com/"},
		UserAgent : "GoScrape",
    		ProxyURL : "http://127.0.0.1:8118"
		Jar : MyCookieJar,
		CallbackF: func (resp *http.Response){
		resp.Body.Close() //VITAL
		}
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
    //Pass options to scraper
    scraper.Scrape(opts)
```
  
# Credits
This module was inspired by:
https://github.com/juliensalinas/go_concurrent_scraping/tree/master/src/go_concurrent_scraper
https://www.admfactory.com/how-to-setup-a-proxy-for-http-client-in-golang/
https://github.com/geziyor/geziyor
https://github.com/headzoo/surf
