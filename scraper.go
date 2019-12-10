package scraper
import(
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)
// STRUCT DEFINITIONS
// Configs defined to pass to Init()
type Config struct {
	// URLSlice ... (Required) A list of strings of the urls of type []string
	URLSlice []string
	// ProxyURL ... (Optional) (e.g. "http://localhost:8118")
	ProxyURL string
	// UserAgent ... (Optional) Set custom user-agent
	UserAgent string
	// Jar ... (Optional) Pass your own cookiejar after signing in for example
	Jar *cookiejar.Jar
	//Callback ... (Optional) I know I just met you, but call me maybe
	CallbackF func(resp *http.Response) interface{}
}
type scraper struct {
	opts *Config
	results []interface{}
}
type ret struct {
	uri string
	result string
}
var userAgent string = "GoScrape"
// Unexported function looped through to grab urls
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
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Printf("Possibly faulty url, make sure of the schema. " + uri)
		chFailedUrls <- uri
		b := &http.Response{}
		chIsFinished <- b
		return
	}
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
// Accepts Config instance, processes it, and returns scraper instance (Doesn't matter if its a copy)
func Init(cfg *Config) (craper scraper) {
	// Stop execution if there are no URLs received
	if len(cfg.URLSlice) == 0 {
		panic("Zero URLs received")
	}
	// If no User-Agent is specified in cfg assign our own
	if cfg.UserAgent != ""{
		userAgent = cfg.UserAgent
	}
	// Create a Callback Function in case none were specified.
	if cfg.CallbackF == nil {
		cfg.CallbackF = func(resp *http.Response) interface{} {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			newStr := buf.String()
			link := resp.Request.URL.String()
			return ret{
				uri: link,
				result: newStr,
			}
		}
		}
	craper = scraper{
		opts: cfg,
	}
	return
}
func (scraper scraper) Results() []interface{} {
	return scraper.results
}
// Starts scraper to loop through supplied URLs using supplied options
func (scraper *scraper) Start () {
	chFailedUrls := make(chan string)
	chIsFinished := make(chan *http.Response)
	for _, uri := range scraper.opts.URLSlice {
		go fetchURL(uri, scraper.opts.ProxyURL, scraper.opts.Jar, chFailedUrls, chIsFinished)
	}
	// Receive messages from every concurrent goroutine. If
	// an url fails, we log it to failedUrls array:
	failedUrls := make([]string, 0)
	for i := 0; i < len(scraper.opts.URLSlice); {
		select {
		case uri := <-chFailedUrls:
			failedUrls = append(failedUrls, uri)
		case resp := <-chIsFinished:
			if resp.Body != nil {
				scraper.results = append(scraper.results, scraper.opts.CallbackF(resp))
				resp.Body.Close()

			}
			i++
		}
	}
	// Print all urls we could not open:
	fmt.Println("Could not fetch these urls: ", failedUrls)
}