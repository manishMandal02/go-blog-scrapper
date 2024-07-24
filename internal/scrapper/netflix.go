package scrapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

func netflix(page *rod.Page, limit int, blog blog) ([]article, error) {
	// netflix blog url
	blogURL := blog.url

	articles := []article{}

	page.MustNavigate(blogURL)

	page.MustWaitDOMStable()

	totalArticlesOnPage := 0

	for totalArticlesOnPage < limit {
		// scroll page to bottom wait 2-5 sec for new articles to be loaded
		page.MustEval(`() => window.scrollTo(0, document.body.scrollHeight)`)

		time.Sleep(time.Second * 3)

		articleEl := page.MustElements("div[data-post-id]")

		totalArticlesOnPage = len(articleEl)
		fmt.Println("🚨 totalArticlesOnPage:", totalArticlesOnPage)
	}

	// ref: skeleton loader selector:- "div.listItemPlaceholder.listItemPlaceholder--withSocialHeader"

	articlesFound, err := getNetflixArticlesOnPage(page, limit, blog)

	if err != nil {
		fmt.Println("error scrapping netflix blog. error:", err)
		return nil, err
	}

	articles = articlesFound

	if len(articles) > limit {
		articles = articles[0:limit]
	}

	return articles, nil
}

func getNetflixArticlesOnPage(page *rod.Page, limit int, blog blog) ([]article, error) {

	articles := []article{}

	articleEl := page.MustElements("div[data-post-id]")

	if len(articleEl) < 1 {
		return nil, errors.New("no articles found")
	}

	// parse extra 10 articles as some are not parsed
	if len(articleEl)+10 > limit {
		articleEl = articleEl[0:limit]
	}

	err := page.GetContext().Err()

	if err != nil {
		fmt.Println("❌ error navigating page.")
	}

	for i, el := range articleEl {

		fmt.Printf("🔂 looping over i:%v, el:%+v", i, el)

		// check if it's a featured article
		hasHeadingContainer, _ := el.Element("div > a > h3")

		fmt.Println("✅ hasHeadingContainer:", hasHeadingContainer)

		article := article{}

		var headingContainer *rod.Element

		if hasHeadingContainer == nil {
			// if article featured, then it has a different heading layout
			headingContainer = el.MustParent().MustElement("div > a:has(h3)")
		} else {
			headingContainer = el.MustElement("div > a:has(h3)")
		}
		fmt.Println("✅ headingContainer:", headingContainer)

		article.url = *headingContainer.MustAttribute("href")

		title, _ := headingContainer.Element("h3 > div")
		desc, _ := headingContainer.Element("h3 ~ div > div")

		if title == nil || desc == nil {
			continue
		}
		article.title = title.MustText()
		article.desc = desc.MustText()

		fmt.Printf("✅ article before date: %+v", article)

		date := ""

		if hasHeadingContainer == nil {
			date = *headingContainer.MustParent().MustElement("time").MustAttribute("datetime")
		} else {
			date = *el.MustElement("time").MustAttribute("datetime")
		}

		// t, err := time.Parse("2006-01-02T15:04-07:00", date)
		t, err := time.Parse(time.RFC3339, date)
		if err != nil {
			fmt.Println("Error parsing time for netflix blog:", err)
			t = time.Now()
		}

		fmt.Println("✅ time:", t)

		article.time = t

		aTag, _ := el.Element("div > a[data-action='open-post']")

		fmt.Println("✅ aTag:", aTag)

		bgImageURL := ""

		if aTag != nil {
			bgImageURL = aTag.MustEval(`() => this.style.backgroundImage.slice(4, -1).replace(/"/g, "")`).Str()
		}

		fmt.Println("✅ bgImageURL:", bgImageURL)

		fmt.Println("🌅 bg image url:", bgImageURL)

		if bgImageURL != "" {
			article.thumbnail = bgImageURL
		} else {
			// use netflix logo as thumbnail if no thumbnail image
			article.thumbnail = blog.logo
		}

		// netflix does not have authors and tags for articles
		article.authors = []string{}
		article.tags = []string{}

		articles = append(articles, article)
	}

	return articles, nil
}
