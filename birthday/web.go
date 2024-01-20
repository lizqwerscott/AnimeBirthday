package birthday

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"sort"
)

type Birthday struct {
	Month int `bson:"month" json:"month"`
	Day   int `bson:"day" json:"day"`
}

type AnimePerson struct {
	Name       string   `bson:"name" json:"name"`
	Url        string   `bson:"url" json:"url"`
	Birthday   Birthday `bson:"birthday" json:"birthday"`
	Reputation int      `bson:"reputation" json:"reputation"`
}

func PrintPersons(x []AnimePerson) {
	log.Printf("len=%d cap=%d\n", len(x), cap(x))

	for i, person := range x {
		log.Printf("person(%d): %s %s %d\n", i, person.Name, person.Url, person.Reputation)
	}
}

func get_birthday_list_from_html(month, day int) ([]AnimePerson, error) {

	get_birthday_url := fmt.Sprintf("https://zh.moegirl.org.cn/Category:%d月%d日", month, day)

	resp, err := http.Get(get_birthday_url)
	if err != nil {
		return nil, errors.Wrapf(err, "get birthday list from html error with month: %d, day: %d", month, day)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(errors.Errorf("status code error: %d %s", resp.StatusCode, resp.Status), "get birthday list from html error with month: %d, day: %d", month, day)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "get birthday list from html error with month: %d, day: %d", month, day)
	}

	persons := make([]AnimePerson, 0)

	not_found_person := make([]string, 0)

	doc.Find(".mw-category-group").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(person_i int, person_li *goquery.Selection) {
			href, err := person_li.Attr("href")
			name := person_li.Text()

			if !err {
				not_found_person = append(not_found_person, name)
			}

			url := fmt.Sprintf("https://zh.moegirl.org.cn%s", href)
			birthday := Birthday{month, day}
			person := AnimePerson{name, url, birthday, 0}
			persons = append(persons, person)
		})
	})

	if len(not_found_person) > 0 {
		return nil, errors.Wrapf(errors.Errorf("not find href with %v", not_found_person), "get birthday list from html error with month: %d, day: %d", month, day)
	}

	return persons, nil
}

func count_page_word(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, errors.Wrapf(err, "count page (%s) error", url)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return -1, errors.Wrapf(errors.Errorf("status code error: %d %s", resp.StatusCode, resp.Status), "count page (%s) error", url)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return -1, errors.Wrapf(err, "count page (%s) error", url)
	}

	return len(body), nil
}

func count_page_word_async(url string, ch chan<- int, err_ch chan<- error) {
	// start := time.Now()

	count, err := count_page_word(url)

	ch <- count
	err_ch <- err

	// end := time.Now()
	// elapse := end.Sub(start)
	// log.Printf("count page Seconds: %f, (%s) \n", elapse.Seconds(), url)
}

func GetAnimePersonBirthdayFromWeb(month, day int) ([]AnimePerson, error) {
	persons, err := get_birthday_list_from_html(month, day)

	if err != nil {
		return nil, errors.Wrapf(err, "Get Anime Person Birthday error with month: %d, day: %d", month, day)
	}

	ch := make(chan int, len(persons))
	err_ch := make(chan error, len(persons))

	for _, person := range persons {
		go count_page_word_async(person.Url, ch, err_ch)
	}

	for i := range persons {
		persons[i].Reputation = <-ch
		err = <-err_ch
		if err != nil {
			return nil, errors.Wrapf(err, "Get Anime Person Birthday error with month: %d, day: %d", month, day)
		}
	}

	sort.SliceStable(persons, func(i, j int) bool {
		return persons[i].Reputation > persons[j].Reputation
	})

	return persons, nil
}
