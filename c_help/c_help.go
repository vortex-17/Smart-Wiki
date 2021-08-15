package c_help

import (
	"fmt"
	"regexp"
	"strings"
)

type Data struct {
	Prev_url   string
	Url        string
	About      string
	Occurences int
	Keyword    string
}

var Num_urls uint64

var Data_list []Data

var Visited map[string]bool

func Filter_link(url string) bool {
	unwanted_list := []string{
		".JPG", ".jpg", "identifier", "Wikipedia", "file", "Special", "Category", "File", "Help", "Template", "Talk", "Module", "Portal",
	}
	for i := range unwanted_list {
		if strings.Contains(url, unwanted_list[i]) {
			return false
		}
	}

	if strings.Contains(url, "https://en.wikipedia.org/wiki/") {
		return true
	}

	return false
}

func Get_words_from(text string) []string {
	words := regexp.MustCompile("\\w+")
	return words.FindAllString(text, -1)
}

func Count_words(words []string, keyword string) int {
	num := 0
	for _, word := range words {
		if word == keyword {
			num++
		}
	}
	return num
}
func Find_word_count(data string, keyword string) int {
	// word_list := Get_words_from(data)
	// num := Count_words(word_list, keyword)
	data = strings.ToLower(data)
	keyword = strings.ToLower(keyword)
	num := strings.Count(data, keyword)
	return num
}

func Sort_data(data_list []Data) []Data {
	for i := 0; i < len(data_list); i++ {
		max := i
		for j := i; j < len(data_list); j++ {
			if data_list[j].Occurences >= data_list[max].Occurences {
				max = j
			}
		}

		data_list[i], data_list[max] = data_list[max], data_list[i]
	}
	return data_list
}

func Print(data_list []Data) {
	for i := 0; i < len(data_list); i++ {
		fmt.Println(i, data_list[i])
	}
}
