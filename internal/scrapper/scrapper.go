package scrapper

import (
	"fmt"
	"os"
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

type Article struct {
	Title     string
	Desc      string
	Time      time.Time
	Tags      []string
	Thumbnail string
	URL       string
}

type blog struct {
	Title    string
	URL      string
	Articles []Article
	Logo     string
}

var Blogs = []blog{
	{
		Title:    "Stripe Engineering",
		URL:      "https://stripe.com/blog/engineering/",
		Logo:     "https://images.ctfassets.net/fzn2n1nzq965/HTTOloNPhisV9P4hlMPNA/cacf1bb88b9fc492dfad34378d844280/Stripe_icon_-_square.svg?q=80&w=1082",
		Articles: []Article{},
	},
	{
		Title:    "Netflix Blog",
		URL:      "https://netflixtechblog.com/",
		Logo:     "https://upload.wikimedia.org/wikipedia/commons/thumb/7/75/Netflix_icon.svg/814px-Netflix_icon.svg.png",
		Articles: []Article{},
	},
	{
		Title:    "Uber  Blog",
		URL:      "https://www.uber.com/en-IN/blog/engineering/",
		Logo:     "https://cdn.icon-icons.com/icons2/2407/PNG/512/uber_icon_146079.png",
		Articles: []Article{},
	},
}

var browser *rod.Browser

var pool rod.Pool[rod.Page]

func StartAll(isHeadless bool) []Article {
	defer utils.FuncExecutionTime()()

	browser = getBrowser(isHeadless)

	defer browser.MustClose()

	// 	concurrently scrape pages
	pool = rod.NewPagePool(3)

	var wg sync.WaitGroup

	allArticles := []Article{}

	for _, blog := range Blogs {
		switch {
		case strings.Contains(blog.Title, "Stripe"):
			wg.Add(1)
			go func() {
				defer wg.Done()
				stripeArticles, err := StripeBlog(-1, isHeadless)
				if err != nil {
					fmt.Println("Error scrapping stripe blog, error:", err)
				} else {
					fmt.Printf("Stripe blog %v\n", len(stripeArticles))
					allArticles = append(allArticles, stripeArticles...)
				}

			}()

		case strings.Contains(blog.Title, "Netflix"):
			wg.Add(1)
			go func() {
				defer wg.Done()
				netflixArticles, err := NetflixBlog(-1, isHeadless)
				if err != nil {
					fmt.Println("Error scrapping netflix blog, error:", err)

				} else {

					fmt.Printf("Netflix blog %v\n", len(netflixArticles))
					allArticles = append(allArticles, netflixArticles...)

				}

			}()

		case strings.Contains(blog.Title, "Uber"):
			wg.Add(1)
			go func() {
				defer wg.Done()
				uberArticles, err := UberBlog(-1, isHeadless)
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
func StripeBlog(limit int, isHeadless bool) ([]Article, error) {
	fmt.Println("Scraping stripe blog, browser:", browser)

	// TODO - also compare browser headless state
	if browser == nil {
		browser = getBrowser(isHeadless)
	}

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := newBrowserPage()

	defer page.Close()

	return stripe(page, limit, Blogs[0])
}

// scrape netflix blog
func NetflixBlog(limit int, isHeadless bool) ([]Article, error) {
	if browser == nil {
		browser = getBrowser(isHeadless)

	}

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := newBrowserPage()

	defer page.Close()

	return netflix(page, limit, Blogs[1])
}

func UberBlog(limit int, isHeadless bool) ([]Article, error) {
	if browser == nil {
		browser = getBrowser(isHeadless)
	}

	limit = utils.SafeMaxLimit(limit, MAX_LIMIT)

	page := newBrowserPage()

	defer page.Close()

	return uber(page, limit, Blogs[2])
}

// scrapper helpers
func getBrowser(isHeadless bool) *rod.Browser {
	// return existing browser if already created
	if browser != nil {
		return browser
	}

	// Set the Chrome path environment variable for Docker
	chromePath := os.Getenv("CHROME_PATH")
	if chromePath == "" {
		chromePath = "/usr/bin/chromium-browser" // default path
	}

	// Initialize a new browser with the specified executable path
	url := launcher.New().Bin(chromePath).Headless(isHeadless).MustLaunch()

	browser = rod.New().ControlURL(url)

	// if !isHeadless {
	// 	browser.ControlURL(launcher.New().Headless(isHeadless).MustLaunch())
	// }

	browser.MustConnect()

	return browser
}

// create pages
// concurrently only when scrapping all Blogs at once
func newBrowserPage() *rod.Page {

	if browser == nil {
		browser = getBrowser(false)
	}

	// if poll is not initialized, then the program does not need concurrent pages
	if pool == nil {
		return browser.MustPage()
	}

	pool = rod.NewPagePool(MAX_PAGE_POOL_LIMIT)

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
