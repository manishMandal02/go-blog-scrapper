package scrapper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/manishmandal02/tech-blog-scrapper/internal/utils"
)

func uber(page *rod.Page, limit int, blog blog) ([]Article, error) {
	// uber blog url
	blogURL := blog.URL

	responses := make(chan []Article)

	articles := []Article{}

	go page.EachEvent(func(ev *proto.PageLoadEventFired) (stop bool) {
		page.MustWaitDOMStable()
		fmt.Println("🌅 uber Page loaded")

		if limit > 0 && len(articles) >= limit {
			// stop if the desired num of articles scrapped
			close(responses)
			return true
		}

		articlesFound, err := getUberArticlesOnPage(page, blog)

		if err != nil {
			fmt.Println("error scrapping uber blog. error:", err)
		}

		fmt.Println("articles scraped:", len(articlesFound))

		nextPageButtons, _ := page.Elements("a[href*='/page/']:has(div > span)")

		var nextPageBtn *rod.Element

		for _, btn := range nextPageButtons {
			nextBtn, _ := btn.Element("div > span")

			if nextBtn == nil {
				continue
			}

			if nextBtn.MustText() != "Next" && nextBtn.MustText() != "View more stories" {
				continue
			}

			nextPageBtn = btn

		}

		responses <- articlesFound

		if nextPageBtn != nil {
			// go to next page
			nextPageBtn.MustClick()
			return false
		}

		// stop
		close(responses)
		return true

	})()

	page.MustNavigate(blogURL)

	for res := range responses {
		fmt.Println("channel res length:", len(res))
		articles = append(articles, res...)
	}

	if len(articles) > limit {
		articles = articles[0:limit]
	}

	return articles, nil
}

func getUberArticlesOnPage(page *rod.Page, blog blog) ([]Article, error) {

	articles := []Article{}

	articleEl := page.MustElements("a[data-baseweb='card']:has(img)")

	if len(articleEl) < 1 {
		return nil, errors.New("no articles found")

	}

	err := page.GetContext().Err()

	if err != nil {
		fmt.Println("❌ error navigating page.")
	}

	for _, el := range articleEl {
		article := Article{}
		article.Title = el.MustElement(" div > div > h5").MustText()
		path := *el.MustAttribute("href")
		article.URL = blog.URL + path

		// desc not found for uber blog
		article.Desc = ""

		date := el.MustElement(" div > div> p").MustText()

		date = strings.Split(date, "/")[0]

		dateSplit := strings.Split(strings.Trim(date, " "), " ")

		day, month, year := dateSplit[0], dateSplit[1], strconv.Itoa(time.Now().Year())

		if len(dateSplit) > 2 {
			year = dateSplit[2]
		}

		t, err := time.Parse("2006-01-02", utils.FormatDateStringUTC(year, month, day))

		if err != nil {
			fmt.Println("Error parsing time for uber blog:", err)
			t = time.Now()
		}
		article.Time = t

		imageTag, _ := el.Element("img")

		if imageTag != nil {
			article.Thumbnail = *imageTag.MustAttribute("src")
		} else {
			// use uber logo as thumbnail if no thumbnail image
			article.Thumbnail = blog.Logo
		}

		tagContainer, _ := el.Element("div > div > div")

		if tagContainer != nil {
			article.Tags = strings.Split(tagContainer.MustText(), ",")
		} else {
			article.Tags = []string{}

		}

		articles = append(articles, article)
	}
	return articles, nil
}
