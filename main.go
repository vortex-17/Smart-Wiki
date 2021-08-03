// This is a tutorial for GoColly/colly
// This will be subsequently be used to scrape wikipedia pages
// For Smart Search project

package main

import (
	"crawler/c_help"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gocolly/colly"
)

var wg sync.WaitGroup

func main() {

	fmt.Println("Welcome to web crawler")
	seed_url := "https://en.wikipedia.org/wiki/Physics"

	// url_channel is the channel used to send links back for finding hyperlinks and data.
	url_channel := make(chan c_help.Data)

	//data_channel is used for sending final data for further operation.
	data_channel := make(chan c_help.Data)

	keyword := "physics"
	wg.Add(3)
	go collect_urls(url_channel, data_channel)
	go collect_urls(url_channel, data_channel)
	go collect_urls(url_channel, data_channel)

	go listen_data(data_channel)

	d := c_help.Data{
		Prev_url:   "Root",
		Url:        seed_url,
		Keyword:    keyword,
		Occurences: 0,
		About:      "",
	}

	crawl(500, d, url_channel, data_channel)

	wg.Wait()
	close(url_channel)

	//print out final data

}

func crawl(n uint64, u c_help.Data, url_channel chan c_help.Data, data_channel chan c_help.Data) {

	if c_help.Num_urls > n {
		//sending end signal
		d := c_help.Data{
			Prev_url:   "",
			Url:        "",
			About:      "",
			Occurences: 0,
			Keyword:    "",
		}
		url_channel <- d
		return
	}

	keyword := u.Keyword
	// fmt.Println("crawler")
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
		colly.MaxDepth(5),
		// colly.Async(true),
	)
	// c.Limit(&colly.LimitRule{
	// 	Delay: 1 * time.Second,
	// })

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println(r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {

		var num int
		num = c_help.Find_word_count(string(r.Body), keyword)
		if num > 1 && c_help.Num_urls <= n {
			url := string(r.Request.URL.String())
			about := strings.Split(url, "/")
			about_word := about[len(about)-1]

			scraped_data := u
			scraped_data.Keyword = keyword
			scraped_data.Occurences = num
			scraped_data.About = about_word
			// fmt.Println(c_help.Num_urls, scraped_data)
			data_channel <- scraped_data
			atomic.AddUint64(&c_help.Num_urls, 1)
		}

	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		check := c_help.Filter_link(e.Request.AbsoluteURL(link))
		if check == true && c_help.Num_urls <= n {
			d := c_help.Data{
				Prev_url:   u.Url,
				Url:        e.Request.AbsoluteURL(link),
				About:      "",
				Occurences: 0,
				Keyword:    keyword,
			}
			url_channel <- d
		} else {
			return
		}

	})

	c.Visit(u.Url)

}

func collect_urls(url_channel chan c_help.Data, data_channel chan c_help.Data) {

	for data := range url_channel {
		// fmt.Println("Crawling URL : ", url)
		if data.Url == "" {
			break
		} else {
			go crawl(500, data, url_channel, data_channel)
		}

	}

	fmt.Println("I am done")
	wg.Done()
}

func listen_data(data_channel chan c_help.Data) {
	for data := range data_channel {
		fmt.Println(data)
	}

}

// func collect_urls(url_channel chan data, data_channel chan data) {
// 	for data := range data_channel {
// 		// fmt.Println("Data Channel : ", data)
// 		url_channel <- data
// 		// fmt.Println("Inserted")
// 	}
// 	fmt.Println("Closing")
// 	close(url_channel)
// }

// func collect_data(url_channel chan data, data_channel chan data) {
// 	for data := range url_channel {
// 		// fmt.Println("URL Channel : ", data)
// 		if data.url == "" {
// 			break
// 		} else {
// 			go crawl(500, data, data_channel)
// 		}

// 		// fmt.Println("Inserted")
// 	}
// 	fmt.Println("Closing Data channel")
// 	close(data_channel)
// }

// func workers(url_channel chan string, crawl_channel chan string){

// }
