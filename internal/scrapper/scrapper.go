package scrapper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type article struct {
	title     string
	desc      string
	time      time.Time
	authors   []string
	tags      []string
	thumbnail string
	url       string
}

type blog struct {
	title    string
	url      string
	articles []article
	logo     string
}

var blogs = []blog{
	{
		title:    "Stripe Engineering",
		url:      "https://stripe.com/blog/engineering",
		logo:     "https://upload.wikimedia.org/wikipedia/commons/b/ba/Stripe_Logo%2C_revised_2016.svg",
		articles: []article{},
	},
	{
		title:    "Netflix",
		url:      "https://netflixtechblog.com/",
		logo:     "https://upload.wikimedia.org/wikipedia/commons/b/ba/Stripe_Logo%2C_revised_2016.svg",
		articles: []article{},
	},
}

func StartAll() {

	for _, blog := range blogs {
		switch {
		case strings.Contains(blog.title, "Stripe"):
			articles, err := stripeBlog(6)
			if err != nil {
				fmt.Println("Error scrapping stripe blog, error:", err)

				return
			}
			fmt.Printf("Stripe blog %+v\n", len(articles))
		case strings.Contains(blog.title, "NULL"):
			articles, err := netflix(12)
			if err != nil {
				fmt.Println("Error scrapping netflix blog, error:", err)
				return
			}
			fmt.Printf("Netflix blog %+v\n", articles)
		default:
			fmt.Println("Unknown blog url.")
		}
	}
}

func netflix(limit int) ([]article, error) {

	// netflix blog url
	blogURL := blogs[1].url
	parsedURL, err := url.Parse(blogURL)

	if err != nil {
		return nil, err
	}

	fmt.Println("parsedURL", parsedURL)

	articles := []article{}

	// scrape

	if len(articles) > limit {
		articles = articles[0:limit]
	}

	return articles, nil
}

func createBrowser(isHeadless bool) *rod.Browser {

	browser := rod.New()

	if !isHeadless {
		browser.ControlURL(launcher.New().Headless(false).MustLaunch())
	}

	return browser
}
