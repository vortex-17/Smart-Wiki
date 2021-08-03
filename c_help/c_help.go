package c_help

import (
	"regexp"
	"strings"
)

type data struct {
	prev_url   string
	url        string
	about      string
	occurences int
	keyword    string
}

var num_urls int

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
	word_list := Get_words_from(data)
	num := Count_words(word_list, keyword)
	return num
}

// func main() {
// 	fmt.Println("this is a helper function")
// }

func Sort_data() {
	return
}
