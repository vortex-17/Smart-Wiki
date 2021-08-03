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

type data struct {
	prev_url   string
	url        string
	about      string
	occurences int
	keyword    string
}

var num_urls uint64

var wg sync.WaitGroup

func main() {

	// var wg sync.WaitGroup
	// wg.Add(2)

	fmt.Println("Welcome to web crawler")
	seed_url := "https://en.wikipedia.org/wiki/Physics"
	url_channel := make(chan data)
	// data_channel := make(chan data)

	keyword := "physics"
	wg.Add(3)
	go collect_urls(url_channel)
	go collect_urls(url_channel)
	go collect_urls(url_channel)
	// go collect_data(url_channel, data_channel)

	// go collect_urls(keyword, url_channel)
	d := data{
		prev_url:   "Root",
		url:        seed_url,
		keyword:    keyword,
		occurences: 0,
		about:      "",
	}
	// url_channel <- seed_url
	// data_channel <- d
	crawl(500, d, url_channel)
	// crawl(seed_url, url_channel)
	wg.Wait()
	close(url_channel)
}

func crawl(n uint64, u data, data_channel chan data) {

	if num_urls > n {
		//sending end signal
		d := data{
			prev_url:   "",
			url:        "",
			about:      "",
			occurences: 0,
			keyword:    "",
		}
		data_channel <- d
		return
	}

	keyword := u.keyword
	// fmt.Println("crawler")
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
		colly.MaxDepth(5),
		// colly.Async(true),
	)
	// infoCollector := c.Clone()

	// c.Limit(&colly.LimitRule{
	// 	Delay: 1 * time.Second,
	// })

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println(r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println(string(r.Body))
		var num int
		// num = find_word_count(string(r.Body), keyword)
		num = c_help.Find_word_count(string(r.Body), keyword)
		if num > 1 && num_urls <= n {
			url := string(r.Request.URL.String())
			about := strings.Split(url, "/")
			about_word := about[len(about)-1]
			// scraped_data := data{
			// 	url:        url,
			// 	about:      about_word,
			// 	occurences: num,
			// 	keyword:    keyword,
			// }

			scraped_data := u
			scraped_data.keyword = keyword
			scraped_data.occurences = num
			scraped_data.about = about_word
			// fmt.Println("Count and link", num)
			fmt.Println(num_urls, scraped_data)
			atomic.AddUint64(&num_urls, 1)
		}

	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// fmt.Println(e.Request.AbsoluteURL(link))
		// check := filter_link(e.Request.AbsoluteURL(link))

		check := c_help.Filter_link(e.Request.AbsoluteURL(link))
		if check == true && num_urls <= n {
			// fmt.Println("link : ", e.Request.AbsoluteURL(link))
			d := data{
				prev_url:   u.url,
				url:        e.Request.AbsoluteURL(link),
				about:      "",
				occurences: 0,
				keyword:    keyword,
			}
			// url_channel <- e.Request.AbsoluteURL(link)
			// fmt.Println(d)
			data_channel <- d

			// c.Visit(e.Request.AbsoluteURL(link))
		} else {
			return
		}

	})

	c.Visit(u.url)

}

func collect_urls(url_channel chan data) {

	for data := range url_channel {
		// fmt.Println("Crawling URL : ", url)
		if data.url == "" {
			break
		} else {
			go crawl(500, data, url_channel)
		}

	}

	fmt.Println("I am done")
	wg.Done()
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
