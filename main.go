// This is a tutorial for GoColly/colly
// This will be subsequently be used to scrape wikipedia pages
// For Smart Search project

package main

import (
	"crawler/c_help"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/gocolly/colly"
)

var wg sync.WaitGroup
var done int
var num int

// var sorted_data []c_help.Data

// var data_map map[string]int

type SafeCounter struct {
	mu sync.Mutex
	v  int
}

type SafeMap struct {
	mu sync.Mutex
	m  map[string]bool
}

func (c *SafeCounter) Inc(n int) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v += n
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	return c.v
}

func (sm *SafeMap) Check(s string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, found := sm.m[s]; !found {
		sm.m[s] = true
		return false
	}

	return true
}

var counter SafeCounter
var checker SafeMap

func SafeClose(ch chan c_help.Data) (justClosed bool) {
	defer func() {
		if recover() != nil {
			// The return result can be altered
			// in a defer function call.
			justClosed = false
		}
	}()

	// assume ch != nil here.
	close(ch)   // panic if ch is closed
	return true // <=> justClosed = true; return
}

func SafeSend(ch chan c_help.Data, value c_help.Data) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()

	ch <- value  // panic if ch is closed
	return false // <=> closed = false; return
}

func Send_data(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "text/html")
	sorted_data := Initialise_crawler(vars["keyword"])
	fmt.Fprintf(w, "<h2> Hello </h2>")
	// fmt.Println(sorted_data)
	for c := range sorted_data {
		fmt.Fprintf(w, "<a href = %s>%s</p>", sorted_data[c].Url, sorted_data[c].About)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h2> Could not find the page</h2>")
}

func main() {
	r := mux.NewRouter()
	s := r.Methods("GET").Subrouter()
	s.NotFoundHandler = http.HandlerFunc(NotFound)
	s.HandleFunc("/data/{keyword}", Send_data)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	// Initialise_crawler()
}

func Initialise_crawler(keyword string) []c_help.Data {

	sorted_data := make([]c_help.Data, 0)
	c_help.Data_list = make([]c_help.Data, 0)

	counter.v = 0
	checker.m = make(map[string]bool)

	keyword = strings.ToLower(keyword)
	split_keywords := strings.Split(keyword, "-")
	keyword = strings.Join(split_keywords, " ")
	var url_keyword string

	fmt.Println(keyword, split_keywords)
	if len(split_keywords) == 1 {
		url_keyword = strings.Title(split_keywords[0])
	} else {
		first_word := strings.Title(split_keywords[0])
		key := strings.Join(split_keywords[1:], "_")
		url_keyword = first_word + "_" + key

	}
	fmt.Println(url_keyword)
	fmt.Println("Welcome to web crawler")
	seed_url := "https://en.wikipedia.org/wiki/"
	seed_url += url_keyword

	fmt.Println(seed_url)

	// url_channel is the channel used to send links back for finding hyperlinks and data.
	url_channel := make(chan c_help.Data)
	fmt.Println("Created data channel")

	done = 1

	wg.Add(3)
	go collect_urls(url_channel)
	go collect_urls(url_channel)
	go collect_urls(url_channel)
	fmt.Println("Started goroutine")
	d := c_help.Data{
		Prev_url:   "Root",
		Url:        seed_url,
		Keyword:    keyword,
		Occurences: 0,
		About:      "",
	}
	url_channel <- d

	fmt.Println("started my service")
	wg.Wait()
	counter.Inc(1)
	fmt.Println("content in channel: ", len(url_channel))

	SafeClose(url_channel)
	fmt.Println("scraping done")
	fmt.Println("Number of links crawler got: ", len(c_help.Data_list))

	//sort according to occurrences
	sorted_data = c_help.Sort_data(c_help.Data_list)
	fmt.Println("data sorted")

	return sorted_data

}

func crawl(n uint64, u c_help.Data, url_channel chan c_help.Data) {
	if uint64(counter.Value()) >= n {
		fmt.Println("Oops, we got all our links")
		return
	}

	keyword := u.Keyword
	// fmt.Println("crawler")
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
		// colly.Async(true),
	)
	// c.Limit(&colly.LimitRule{
	// 	Delay: 1 * time.Second,
	// })

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting :", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
		var num int
		num = c_help.Find_word_count(string(r.Body), keyword)
		if num > 1 && uint64(counter.Value()) < n {
			url := string(r.Request.URL.String())
			about := strings.Split(url, "/")
			about_word := about[len(about)-1]

			scraped_data := u
			scraped_data.Keyword = keyword
			scraped_data.Occurences = num
			scraped_data.About = about_word
			// fmt.Println(c_help.Num_urls, scraped_data)
			counter.Inc(1)
			fmt.Println("data: ", uint64(counter.Value()), scraped_data)
			c_help.Data_list = append(c_help.Data_list, scraped_data)
		}

	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err, counter.Value())
		if counter.Value() == 0 {
			log.Println("Root URL is bad")
			close(url_channel)
		}
		return
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if uint64(counter.Value()) <= n {
			link := e.Attr("href")
			check := c_help.Filter_link(e.Request.AbsoluteURL(link))
			if check == true && uint64(counter.Value()) <= n && checker.Check(e.Request.AbsoluteURL(link)) == false {
				d := c_help.Data{
					Prev_url:   u.Url,
					Url:        e.Request.AbsoluteURL(link),
					About:      "",
					Occurences: 0,
					Keyword:    keyword,
				}

				if SafeSend(url_channel, d) == true {
					fmt.Println("Oops, not a safe send")
				}
			}
		} else {
			return
		}

	})

	c.Visit(u.Url)
	return

}

func collect_urls(url_channel chan c_help.Data) {

	for data := range url_channel {

		// fmt.Println(counter.Value(), data.Url)
		if counter.Value() >= 100 {
			fmt.Println("number of url >>>> 100")
			break
		} else {
			// fmt.Println("Crawling URL : ", data.Url)
			go crawl(100, data, url_channel)
		}

	}

	fmt.Println("I am done")
	wg.Done()
}
