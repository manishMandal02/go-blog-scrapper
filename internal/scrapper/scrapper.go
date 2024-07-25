package scrapper

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/manishmandal02/tech-blog-scrapper/internal/utils"
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
		title:    "Netflix Blog",
		url:      "https://netflixtechblog.com/",
		logo:     "https://upload.wikimedia.org/wikipedia/commons/thumb/7/75/Netflix_icon.svg/814px-Netflix_icon.svg.png",
		articles: []article{},
	},
	{
		title:    "Uber  Blog",
		url:      "https://www.uber.com/en-IN/blog/engineering/",
		logo:     "https://upload.wikimedia.org/wikipedia/commons/thumb/5/58/Uber_logo_2018.svg/1600px-Uber_logo_2018.svg.png?20180914002846",
		articles: []article{},
	},
}

var browser *rod.Browser

func StartAll() {

	browser = getBrowser(false)

	for _, blog := range blogs {
		switch {
		case strings.Contains(blog.title, "Stripe"):
			stripeArticles, err := StripeBlog(200)
			if err != nil {
				fmt.Println("Error scrapping stripe blog, error:", err)
				return
			}
			fmt.Printf("Stripe blog %v\n", len(stripeArticles))
		case strings.Contains(blog.title, "Netflix"):
			netflixArticles, err := NetflixBlog(120)
			if err != nil {
				fmt.Println("Error scrapping netflix blog, error:", err)
				return
			}
			fmt.Printf("Netflix blog %v\n", len(netflixArticles))
		case strings.Contains(blog.title, "Uber"):
			uberArticles, err := UberBlog(120)
			if err != nil {
				fmt.Println("Error scrapping uber blog, error:", err)
				return
			}
			fmt.Printf("Uber blog %v\n", len(uberArticles))
		default:
			fmt.Println("Unknown blog url.")
		}
	}
}

// scrape stripe blog
func StripeBlog(limit int) ([]article, error) {
	if browser == nil {
		browser = getBrowser(true)
	}

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := browser.MustPage()

	return stripe(page, limit, blogs[0])
}

// scrape netflix blog
func NetflixBlog(limit int) ([]article, error) {
	if browser == nil {
		browser = getBrowser(true)
	}

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := browser.MustPage()

	return netflix(page, limit, blogs[1])
}

func UberBlog(limit int) ([]article, error) {
	browser = getBrowser(false)

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := browser.MustPage()

	return uber(page, limit, blogs[2])
}

func getBrowser(isHeadless bool) *rod.Browser {
	// return existing browser if already created
	if browser != nil {
		return browser
	}

	browser = rod.New()

	if !isHeadless {
		browser.ControlURL(launcher.New().Headless(false).MustLaunch())
	}

	browser.MustConnect()

	return browser
}
