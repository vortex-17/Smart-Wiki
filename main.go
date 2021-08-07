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
var done int
var num int

// var data_map map[string]int

func main() {

	fmt.Println("Welcome to web crawler")
	seed_url := "https://en.wikipedia.org/wiki/Physics"

	// url_channel is the channel used to send links back for finding hyperlinks and data.
	url_channel := make(chan c_help.Data)

	//data_channel is used for sending final data for further operation.
	// data_channel := make(chan c_help.Data)

	// data_map := make(map[string]int)
	done = 1

	keyword := "physics"
	wg.Add(3)
	go collect_urls(url_channel)
	go collect_urls(url_channel)
	go collect_urls(url_channel)
	// go collect_urls(url_channel, data_channel)

	// go listen_data(data_channel)

	d := c_help.Data{
		Prev_url:   "Root",
		Url:        seed_url,
		Keyword:    keyword,
		Occurences: 0,
		About:      "",
	}

	// c_help.Visited = make(map[string]bool)
	// c_help.Visited[seed_url] = false
	go crawl(100, d, url_channel)

	wg.Wait()
	// close(url_channel)
	// close(data_channel)
	if _, ok := <-url_channel; !ok {
		close(url_channel)
	}
	fmt.Println("scraping done")
	fmt.Println("Number of links crawler got: ", len(c_help.Data_list))

	//sort according to occurrences
	sorted_data := c_help.Sort_data(c_help.Data_list)
	fmt.Println("data sorted")

	//print out final data
	c_help.Print(sorted_data)
	// fmt.Println(len(sorted_data))
	// fmt.Scanln()

}

func crawl(n uint64, u c_help.Data, url_channel chan c_help.Data) {

	// if done == 0 {
	// 	return
	// }

	if c_help.Num_urls >= n {
		//sending end signal
		// fmt.Println(c_help.Num_urls)
		d := c_help.Data{
			Prev_url:   "",
			Url:        "",
			About:      "",
			Occurences: 0,
			Keyword:    "",
		}

		// if _, ok := <-url_channel; !ok {
		// 	return
		// }
		url_channel <- d
		// close(url_channel)
		return
	}

	keyword := u.Keyword
	// fmt.Println("crawler")
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
		colly.MaxDepth(5),
		colly.Async(true),
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
			atomic.AddUint64(&c_help.Num_urls, 1)
			c_help.Data_list = append(c_help.Data_list, scraped_data)
			// data_channel <- scraped_data
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
			// c_help.Visited[e.Request.AbsoluteURL(link)] = false
			url_channel <- d
		} else {
			return
		}

	})

	// if c_help.Visited[u.Url] == false {
	// 	c_help.Visited[u.Url] = true
	// 	c.Visit(u.Url)
	// }

	c.Visit(u.Url)

}

func collect_urls(url_channel chan c_help.Data) {

	for data := range url_channel {
		// fmt.Println("Crawling URL : ", url)
		if data.Url == "" {
			done = 0
			break
		}
		if done == 1 {
			go crawl(100, data, url_channel)
		}

	}

	fmt.Println("I am done")
	wg.Done()
}

func listen_data(data_channel chan c_help.Data) {
	for data := range data_channel {
		// fmt.Println(data)
		c_help.Data_list = append(c_help.Data_list, data)
		num++
	}

	// close(data_channel)

}
