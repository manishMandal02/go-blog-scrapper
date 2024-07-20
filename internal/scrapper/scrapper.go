package scrapper

import (
	"fmt"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type article struct {
	title     string
	desc      string
	time      time.Time
	author    []string
	tags      []string
	thumbnail string
	url       string
}

type blog struct {
	title    string
	url      string
	articles []article
}

var blogs = []blog{
	{
		title:    "Stripe engineering",
		url:      "https://stripe.com/blog/engineering",
		articles: []article{},
	},
}

func Start() {

	for _, blog := range blogs {
		switch blog.title {
		case "Stripe engineering":
			articles, err := stripeBlog()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%+v\n", articles)

		default:

		}

	}
}

func stripeBlog() ([]article, error) {
	// strip blog url
	blogURL := blogs[0].url

	parsedURL, err := url.Parse(blogURL)

	if err != nil {
		return nil, err
	}

	articles := []article{}

	//  initialize colly with only stripe url
	c := colly.NewCollector(colly.AllowedDomains(parsedURL.Hostname()))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("article", func(el *colly.HTMLElement) {
		article := article{}
		article.title = el.ChildText("h1 > a")
		article.url = el.ChildAttr("h1 > a", "href")
		article.desc = el.ChildText("div.BlogIndexPost__body > p")
		date := el.ChildAttr("time", "datetime")
		// format date to to RFC3339 layout
		// lastDash := strings.LastIndex(date, "-")
		// date = date[:lastDash] + ":00Z" + date[lastDash+1:]
		t, err := time.Parse("2006-01-02T15:04-07:00", date)
		if err != nil {
			fmt.Println("Error parsing time for stripe blog:", err)
		}
		article.time = t
		article.thumbnail = el.ChildAttr("picture > img", "src")
		tag := ""
		article.author = el.DOM.Find("figure ").Map(func(i int, s *goquery.Selection) string {
			// get tag (using author department as tag)
			tag = s.Find("figcaption > span").Text()

			// returns author name
			return s.Find("figcaption > a").Text()
		})
		article.tags = []string{tag}

		articles = append(articles, article)

	})

	c.OnHTML("a.BlogCategoryPagination__directionLink", func(el *colly.HTMLElement) {
		nextPage := el.Attr("href")

		fmt.Println("Visiting next page:", nextPage)
		c.Visit("https://" + parsedURL.Hostname() + nextPage)
	})

	c.Visit(blogURL)

	return articles, err
}
