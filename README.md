# Description
GoScrape is a simple web scraper written in Go. It supports concurrency, proxy, user-agent, and cookies.

# Installation
```go get github.com/mohd-taba/GoScrape```

# Usage
```golang
import "github.com/mohd-taba/GoScrape"

// Define your configurations
	cfg := scraper.Config{
	URLSlice:  []string {"http://duckduckgo.com", "http://api.myip.com"},
	ProxyURL:  "http://127.0.0.1:8118",
	UserAgent: "TESTIS",
	Jar:       myJar,
}

// Create your scraper
scraperInstance := scraper.Init(cfg)
// Start scraping
scraperInstance.Start()
```

# Callback Function
The callback function is a function supplied to manipulate the response after a request has been made, if none were supplied a default function will be defined.

### Example

```golang
/*Import and header...
The defined function must accept a *http.Response and return an interface{}
The return type was specified as interface to allow for extra degrees of freedom of code*/

func cbf(r *http.Response) interface{}{
	return "Frick" //You can return almost any type
}

cfg := scraper.Config{
	URLSlice:  []string {"http://duckduckgo.com", "http://api.myip.com"},
	ProxyURL:  "http://127.0.0.1:8118",
	UserAgent: "TESTIS",
	CallbackF: cbf,
	Jar:       myJar,
}

```

# Cookies
```golang
import (
	"net/http/cookiejar"
	"github.com/mohd-taba/GoScrape"
	)
	// Define jar options if necessary
jarOptions := cookiejar.Options{
        PublicSuffixList: publicsuffix.List,
    }
	
	// Create a new jar
    jar, err := cookiejar.New(&jarOptions)
    if err != nil {
        log.Fatal(err)
    }
	
    // Do stuff using the cookie jar (e.g, sign-in)
    client := http.Client{Jar: jar}
    resp, err := client.Get("https://cookiesite.com/login.php?=username&password")
    if err != nil {
        log.Fatal(err)
    }
	
	//Create a scraper.Config instance with the same jar specified
	cfg := scraper.Config{
		URLSlice: []string {"http://website.com", "https://example.com"}
		Jar: jar
	}
	
	//Create Scraper
	scraper.Init(cfg)
	scraper.Start()
	
