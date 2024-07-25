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

func uber(page *rod.Page, limit int, blog blog) ([]article, error) {
	// uber blog url
	blogURL := blog.url

	responses := make(chan []article)

	articles := []article{}

	go page.EachEvent(func(ev *proto.PageLoadEventFired) (stop bool) {
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

			fmt.Println("🚨 nextBtn text:", nextBtn.MustText() == "View more stories")
			if nextBtn.MustText() != "Next" && nextBtn.MustText() != "View more stories" {
				continue
			}

			nextPageBtn = btn

		}

		responses <- articlesFound

		if nextPageBtn == nil {
			// stop
			close(responses)
			return true
		}

		// go to next page
		nextPageBtn.MustClick()
		return false
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

func getUberArticlesOnPage(page *rod.Page, blog blog) ([]article, error) {

	articles := []article{}

	articleEl := page.MustElements("a[data-baseweb='card']:has(img)")

	if len(articleEl) < 1 {
		return nil, errors.New("no articles found")

	}

	err := page.GetContext().Err()

	if err != nil {
		fmt.Println("❌ error navigating page.")
	}

	for _, el := range articleEl {
		article := article{}
		article.title = el.MustElement(" div > div > h5").MustText()
		article.url = *el.MustAttribute("href")
		// desc not found for uber blog
		article.desc = ""

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
		article.time = t

		imageTag, _ := el.Element("img")

		if imageTag != nil {
			article.thumbnail = *imageTag.MustAttribute("src")
		} else {
			// use uber logo as thumbnail if no thumbnail image
			article.thumbnail = blog.logo
		}

		article.authors = []string{}

		tag := el.MustElement("div > div > div").MustText()

		article.tags = strings.Split(tag, ",")

		articles = append(articles, article)
	}
	return articles, nil
}
