package scrapper

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/manishmandal02/tech-blog-scrapper/internal/utils"
)

const MAX_LIMIT int = 100

// max concurrent pages to be opened
const MAX_PAGE_POOL_LIMIT int = 3

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

var pool rod.Pool[rod.Page]

func StartAll() []article {
	defer utils.FuncExecutionTime()()

	browser = getBrowser(false)

	defer browser.MustClose()

	// 	concurrently scrape pages
	pool = rod.NewPagePool(3)

	var wg sync.WaitGroup

	allArticles := []article{}

	for _, blog := range blogs {
		switch {
		case strings.Contains(blog.title, "Stripe"):
			wg.Add(1)
			go func() {
				defer wg.Done()
				stripeArticles, err := StripeBlog(200)
				if err != nil {
					fmt.Println("Error scrapping stripe blog, error:", err)
				} else {
					fmt.Printf("Stripe blog %v\n", len(stripeArticles))
					allArticles = append(allArticles, stripeArticles...)
				}

			}()

		case strings.Contains(blog.title, "Netflix"):
			wg.Add(1)
			go func() {
				defer wg.Done()
				netflixArticles, err := NetflixBlog(120)
				if err != nil {
					fmt.Println("Error scrapping netflix blog, error:", err)

				} else {

					fmt.Printf("Netflix blog %v\n", len(netflixArticles))
					allArticles = append(allArticles, netflixArticles...)

				}

			}()

		case strings.Contains(blog.title, "Uber"):
			wg.Add(1)
			go func() {
				defer wg.Done()
				uberArticles, err := UberBlog(120)
				if err != nil {
					fmt.Println("Error scrapping uber blog, error:", err)
				} else {
					fmt.Printf("Uber blog %v\n", len(uberArticles))
					allArticles = append(allArticles, uberArticles...)
				}

			}()

		default:
			fmt.Println("Unknown blog url.")
		}
	}

	wg.Wait()

	return allArticles
}

// scrape stripe blog
func StripeBlog(limit int) ([]article, error) {
	if browser == nil {
		browser = getBrowser(true)
	}

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := newBrowserPage()

	defer page.Close()

	return stripe(page, limit, blogs[0])
}

// scrape netflix blog
func NetflixBlog(limit int) ([]article, error) {
	if browser == nil {
		browser = getBrowser(true)
	}

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := newBrowserPage()

	defer page.Close()

	return netflix(page, limit, blogs[1])
}

func UberBlog(limit int) ([]article, error) {
	browser = getBrowser(false)

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := newBrowserPage()

	defer page.Close()

	return uber(page, limit, blogs[2])
}

// scrapper helpers

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

// create pages
// concurrently only when scrapping all blogs at once
func newBrowserPage() *rod.Page {

	if browser == nil {
		browser = getBrowser(false)
	}

	// if poll is not initialized, then the program does not need concurrent pages
	if pool == nil {
		pool = rod.NewPagePool(MAX_PAGE_POOL_LIMIT)
		return browser.MustPage()
	}

	createPage := func() (*rod.Page, error) {
		return browser.Page(proto.TargetCreateTarget{})
	}

	page, err := pool.Get(createPage)

	if err != nil {
		fmt.Println("Error getting page from pool, error:", err)
		return nil
	}

	return page
}
