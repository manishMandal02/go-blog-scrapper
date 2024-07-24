package scrapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func stripe(page *rod.Page, limit int, blog blog) ([]article, error) {
	// strip blog url
	blogURL := blog.url

	responses := make(chan []article)

	articles := []article{}

	go page.EachEvent(func(ev *proto.PageLoadEventFired) (stop bool) {
		fmt.Println("🌅 Page loaded")

		if limit > 0 && len(articles) >= limit {
			// stop if the desired num of articles scrapped
			close(responses)
			return true
		}

		articlesFound, err := getStripeArticlesOnPage(page, blog)

		if err != nil {
			fmt.Println("error scrapping stripe blog. error:", err)
		}

		fmt.Println("articles scraped:", len(articlesFound))

		nextPageButtons, _ := page.Elements("a.BlogCategoryPagination__directionLink")

		var nextPageBtn *rod.Element

		for _, btn := range nextPageButtons {
			if text, _ := btn.Text(); text != "Next" {
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

		fmt.Println("🚨 Total articles:", len(articles))

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

func getStripeArticlesOnPage(page *rod.Page, blog blog) ([]article, error) {

	articles := []article{}

	articleEl := page.MustElements("article")

	if len(articleEl) < 1 {
		return nil, errors.New("no articles found")

	}

	err := page.GetContext().Err()

	if err != nil {
		fmt.Println("❌ error navigating page.")
	}

	for _, el := range articleEl {
		article := article{}
		article.title = el.MustElement(" h1 > a").MustText()
		article.url = *el.MustElement(" h1 > a").MustAttribute("href")
		article.desc = el.MustElement("div.BlogIndexPost__body > p").MustText()

		date := *el.MustElement("time").MustAttribute("datetime")
		t, err := time.Parse("2006-01-02T15:04-07:00", date)
		if err != nil {
			fmt.Println("Error parsing time for stripe blog:", err)
			t = time.Now()
		}
		article.time = t

		imageTag, _ := el.Element("picture > img")

		if imageTag != nil {
			article.thumbnail = *imageTag.MustAttribute("src")
		} else {
			// use stripe logo as thumbnail if no thumbnail image
			article.thumbnail = blog.logo
		}

		authorContainer := el.MustElement("div.BlogIndexPost__authorList").MustElements("figure")

		for _, authorNode := range authorContainer {

			author := authorNode.MustElement("figcaption > a").MustText()
			article.authors = append(article.authors, author)
			// tags
			tag := authorNode.MustElement("figcaption > span").MustText()
			article.tags = append(article.tags, tag)

		}

		articles = append(articles, article)
	}
	return articles, nil
}