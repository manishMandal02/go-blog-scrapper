package scrapper

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/manishmandal02/tech-blog-scrapper/internal/utils"
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
		title:    "Stripe Engineering",
		url:      "https://stripe.com/blog/engineering",
		articles: []article{},
	},
	{
		title:    "Netflix",
		url:      "https://netflixtechblog.com/",
		articles: []article{},
	},
}

func StartAll() {
	for _, blog := range blogs {
		switch {
		case strings.Contains(blog.title, "Uber"):
			articles, err := stripeBlog(2)
			if err != nil {
				fmt.Println("Error scrapping stripe blog, error:", err)

				return
			}
			fmt.Printf("Stripe blog %+v\n", len(articles))
		case strings.Contains(blog.title, "Netflix"):
			articles, err := uber(12)
			if err != nil {
				fmt.Println("Error scrapping uber blog, error:", err)
				return
			}
			fmt.Printf("Uber blog %+v\n", articles)
		default:
			fmt.Println("Unknown blog url.")
		}
	}
}

func uber(limit int) ([]article, error) {

	// uber blog url
	blogURL := blogs[1].url
	parsedURL, err := url.Parse(blogURL)

	if err != nil {
		return nil, err
	}

	fmt.Println("parsedURL", parsedURL)

	articles := []article{}

	//  initialize colly with only uber blog
	c := colly.NewCollector(colly.AllowedDomains(parsedURL.Hostname(), "medium.com"))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("colly error res: %+v\n", r.Request)
		fmt.Println("uber-blog: Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("a[data-baseweb='card']", func(el *colly.HTMLElement) {
		fmt.Println("HTML loaded, article card el:", el)
		article := article{}
		article.title = el.ChildText("div > div > h5")
		article.url = el.Attr("href")
		// no desc on uber blog
		article.desc = ""

		// date format: 18 July / Global
		dateRaw := el.ChildText("div > div> p")
		dateParsed := strings.Split(dateRaw, "/")[0]

		// [14, December, 2024] *year isn't alway present
		dateFragments := strings.Split(dateParsed, " ")

		day, month, year := dateFragments[0], dateFragments[1], time.Now().Year()

		if len(dateFragments) > 1 && dateFragments[2] != "" {
			y, err := strconv.Atoi(dateFragments[2])
			if err == nil {
				year = y
			}
		}

		formattedDate := utils.FormatDateStringUTC(string(year), month, day)

		t, err := time.Parse("2006-01-02", formattedDate)
		if err != nil {
			fmt.Println("uber-blog: Error parsing time. error:", err)
			t = time.Now()
		}
		article.time = t
		article.thumbnail = el.ChildAttr("img", "src")
		article.author = []string{}

		tags := el.ChildText("div > div > div")

		article.tags = strings.Split(tags, ",")

		articles = append(articles, article)
	})

	// navigate pages
	c.OnHTML("a[href*='/page/']", func(el *colly.HTMLElement) {
		hasNextPage := el.ChildText("div > span")

		if hasNextPage != "Next" {
			return
		}

		if limit > 0 && len(articles) >= limit {
			return
		}

		fmt.Println("uber-blog: visiting next page...")
		c.Visit(el.Attr("href"))
	})

	c.Visit(blogURL)

	if len(articles) > limit {
		articles = articles[0:limit]
	}

	return articles, nil
}

func stripeBlog(limit int) ([]article, error) {
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

		t, err := time.Parse("2006-01-02T15:04-07:00", date)
		if err != nil {
			fmt.Println("Error parsing time for stripe blog:", err)
			t = time.Now()
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

	// navigate pages
	c.OnHTML("a.BlogCategoryPagination__directionLink", func(el *colly.HTMLElement) {
		nextPage := el.Attr("href")

		if limit > 0 && len(articles) >= limit {
			return
		}

		fmt.Println("Visiting next page:", nextPage)
		c.Visit("https://" + parsedURL.Hostname() + nextPage)
	})

	c.Visit(blogURL)

	if len(articles) > limit {
		articles = articles[0:limit]
	}

	return articles, err
}
