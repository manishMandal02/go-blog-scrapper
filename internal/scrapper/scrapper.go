package scrapper

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const MAX_LIMIT int = 100

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
		logo:     "https://upload.wikimedia.org/wikipedia/commons/thumb/7/75/Netflix_icon.svg/814px-Netflix_icon.svg.png",
		articles: []article{},
	},
}

var browser *rod.Browser

func StartAll() {

	browser = createBrowser(false)

	for _, blog := range blogs {
		switch {
		case strings.Contains(blog.title, "NULL"):
			articles, err := StripeBlog(6)
			if err != nil {
				fmt.Println("Error scrapping stripe blog, error:", err)

				return
			}
			fmt.Printf("Stripe blog %v\n", len(articles))
		case strings.Contains(blog.title, "Netflix"):
			articles, err := NetflixBlog(-1)
			if err != nil {
				fmt.Println("Error scrapping netflix blog, error:", err)
				return
			}
			fmt.Printf("Netflix blog %v\n", len(articles))
		default:
			fmt.Println("Unknown blog url.")
		}
	}
}

// scrape stripe blog
func StripeBlog(limit int) ([]article, error) {
	if browser == nil {
		browser = createBrowser(true)
	}

	// max limit to 100 articles
	if limit > MAX_LIMIT || limit < 1 {
		limit = MAX_LIMIT
	}

	page := browser.MustConnect().MustPage()

	return stripe(page, limit, blogs[0])
}

// scrape netflix blog
func NetflixBlog(limit int) ([]article, error) {
	if browser == nil {
		browser = createBrowser(true)
	}

	page := browser.MustConnect().MustPage()

	// max limit to 100 articles
	if limit > MAX_LIMIT || limit < 1 {
		limit = MAX_LIMIT
	}

	return netflix(page, limit, blogs[1])
}

func createBrowser(isHeadless bool) *rod.Browser {

	browser := rod.New()

	if !isHeadless {
		browser.ControlURL(launcher.New().Headless(false).MustLaunch())
	}

	return browser
}
