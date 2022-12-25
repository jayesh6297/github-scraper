package githubscrape

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// scrapes data of repos based on orgnizations name
func ScrapeUser(userID string) {
	total := TotalPages(userID)
	// make chan of total available repos
	// 1 page has 30 repos
	reposUserCh := make(chan *Repo, total*30)
	var wg sync.WaitGroup

	wg.Add(total)
	for p := 1; p <= total; p++ {
		// generate url per page
		url := fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d", userID, p)
		ProcessPageUser(url, userID, reposUserCh, &wg)
	}

	wg.Wait()
	close(reposUserCh)

	// create file after closing channel
	CreateFileUser(fmt.Sprintf("users/%s.json", userID), reposUserCh)
}

// get total pages of repositories
func TotalPagesUser(userID string) int {
	rootURL := fmt.Sprintf("https://github.com/%s?tab=repositories", userID)
	res, err := http.Get(rootURL)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}

	pages, ok := doc.Find("#user-repositories-list").Find("div.pagination").Find("em.current").Attr("data-total-pages")
	if ok {
		pagesInt, _ := strconv.Atoi(pages)
		return pagesInt
	}
	return 1
}

// concurrently process per page
func ProcessPageUser(url string, orgName string, ch chan<- *Repo, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}

	selection := doc.Find("#user-repositories-list").Find("ul").Find("li")
	selection.Each(ProcessRepoUser(orgName, ch))
}

// process repo data for single repo found
func ProcessRepoUser(userID string, ch chan<- *Repo) func(i int, s *goquery.Selection) {
	return func(i int, s *goquery.Selection) {
		baseName := s.Find("a[itemprop='name codeRepository']")
		title := ClearString(baseName.Text())
		link, _ := baseName.Attr("href")
		description := ClearString(s.Find("p[itemprop='description']").Text())
		language := s.Find("span[itemprop='programmingLanguage']").Text()
		forks := ClearString(s.Find(fmt.Sprintf("a[href='/%s/%s/network/members']", userID, title)).Text())
		stars := ClearString(s.Find(fmt.Sprintf("a[href='/%s/%s/stargazers']", userID, title)).Text())
		ch <- &Repo{
			Title:       title,
			Link:        URL + link,
			Description: description,
			Language:    language,
			Forks:       forks,
			Stars:       stars,
		}
	}
}

// creates json file for per org
func CreateFileUser(path string, ch <-chan *Repo) {
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	var ll []*Repo
	for i := range ch {
		ll = append(ll, i)
	}
	b, _ := json.MarshalIndent(map[string]any{"count": len(ll), "repos": ll}, "", "  ")
	f.Write(b)
}

// cleans the string
func ClearStringUser(s string) string {
	return strings.Trim(
		strings.ReplaceAll(
			s,
			"\n",
			"",
		),
		" ",
	)
}
