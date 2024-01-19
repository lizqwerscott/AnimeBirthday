package birthday

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

type Birthday struct {
	Month int `json:"month"`
	Day   int `json:"day"`
}

type AnimePerson struct {
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	Birthday   Birthday `json:"birthday"`
	Reputation int      `json:"reputation"`
}

func PrintPersons(x []AnimePerson) {
	fmt.Printf("len=%d cap=%d\n", len(x), cap(x))

	for i, person := range x {
		fmt.Printf("person(%d): %s %s %d\n", i, person.Name, person.Url, person.Reputation)
	}
}

func get_birthday_list_from_html(month, day int) (value []AnimePerson) {

	get_birthday_url := fmt.Sprintf("https://zh.moegirl.org.cn/Category:%d月%d日", month, day)

	resp, err := http.Get(get_birthday_url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var persons []AnimePerson

	doc.Find(".mw-category-group").Each(func(i int, s *goquery.Selection) {
		// title := s.Find("h3").Text()
		// fmt.Printf("category group(%d): %s\n", i, title)

		s.Find("a").Each(func(person_i int, person_li *goquery.Selection) {
			href, err := person_li.Attr("href")
			name := person_li.Text()

			if err {
				url := fmt.Sprintf("https://zh.moegirl.org.cn%s", href)

				// fmt.Printf("person (%d) (%s): %s\n", i, name, url)
				birthday := Birthday{month, day}
				person := AnimePerson{name, url, birthday, 0}
				persons = append(persons, person)
			}
		})
	})

	return persons
}

func count_page_word(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return -1
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("count page (%s) status code error: %d %s", url, resp.StatusCode, resp.Status)
		return -1
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return -1
	}

	return len(body)
}

func count_page_word_async(url string, ch chan<- int) {
	start := time.Now()

	count := count_page_word(url)

	ch <- count

	end := time.Now()
	elapse := end.Sub(start)
	fmt.Printf("count page Seconds: %f, (%s) \n", elapse.Seconds(), url)
}

func GetAnimePersonBirthday(month, day int) []AnimePerson {
	persons := get_birthday_list_from_html(month, day)

	ch := make(chan int, len(persons))

	for _, person := range persons {
		go count_page_word_async(person.Url, ch)
	}

	for i := range persons {
		persons[i].Reputation = <-ch
	}

	sort.SliceStable(persons, func(i, j int) bool {
		return persons[i].Reputation > persons[j].Reputation
	})

	return persons
}
