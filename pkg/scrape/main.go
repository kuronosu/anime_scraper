package scrape

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/kuronosu/anime_scraper/pkg/config"
)

func ScrapeDetails(
	schema *config.PageSchema,
	urls []string,
	detailsCollector DetailsCollector,
	async bool,
) map[string]string {
	c := colly.NewCollector()
	if async {
		c = colly.NewCollector(colly.Async(async))
	}
	if schema.Cloudflare {
		c.WithTransport(GetCloudFlareRoundTripper())
	}
	mutex := &sync.RWMutex{}

	errorUrls := make(map[string]string)
	responseCounter := 0
	total := len(urls)

	for _, field := range schema.Detail.Fields {
		func(c *colly.Collector, field config.Field) {
			c.OnHTML(field.Selector, func(e *colly.HTMLElement) {
				tmp := field.SafeCompile(e)
				if field.IsString() && field.Contains != nil {
					if field.Contains.Raw {
						if !strings.Contains(e.Text, field.Contains.String) {
							return
						}
					} else {
						if !strings.Contains(fmt.Sprint(tmp), field.Contains.String) {
							return
						}
					}
				}
				detailsCollector.CollectField(e.Request.URL.String(), field.Name, tmp)
			})
		}(c, field)
	}

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		responseCounter++
		fmt.Print("\r[Progress] ", responseCounter, "/", total, " [Errors] ", len(errorUrls))
	})

	c.OnError(func(r *colly.Response, err error) {
		mutex.Lock()
		errorUrls[r.Request.URL.String()] = err.Error()
		mutex.Unlock()
		responseCounter++
		fmt.Print("\r[Progress] ", responseCounter, "/", total, " [Errors] ", len(errorUrls))
	})

	// fmt.Print("\r[Progress] ", counter, "/", total)
	requestCounter := 0
	for _, url := range urls {
		c.Visit(url)
		requestCounter++
		if requestCounter >= 500 {
			requestCounter = 0
			c.Wait()
		}
	}
	c.Wait()
	time.Sleep(300 * time.Millisecond)
	fmt.Println(" Done")
	return errorUrls
}

func ScrapeList(schema *config.PageSchema, url string) config.ParsedLinks {
	c := colly.NewCollector()
	if schema.Cloudflare {
		c.WithTransport(GetCloudFlareRoundTripper())
	}
	results := make(config.ParsedLinks)

	visitedUrls := make(map[string]bool)

	c.OnHTML(schema.List.ContainerSelector, func(e *colly.HTMLElement) {
		results.Extend(config.ParseLinkData(schema.List.SafeCompile(e)))
	})

	pagPrevCount := 0
	pagNextCount := 0

	c.OnHTML(schema.List.Pagination.Next.Selector, func(e *colly.HTMLElement) {
		// fmt.Println(schema.List.Pagination.Next.GetValue(e))
		visitedUrls[e.Request.URL.String()] = true
		pageUrl := schema.List.Pagination.Next.GetValue(e)
		if _, ok := visitedUrls[pageUrl]; !ok {
			pagNextCount++
			if schema.List.Pagination.Next.Limit == -1 || pagNextCount <= schema.List.Pagination.Next.Limit {
				e.Request.Visit(pageUrl)
			}
		}
	})

	c.OnHTML(schema.List.Pagination.Previous.Selector, func(e *colly.HTMLElement) {
		// fmt.Println(schema.List.Pagination.Previous.GetValue(e))
		visitedUrls[e.Request.URL.String()] = true
		pageUrl := schema.List.Pagination.Previous.GetValue(e)
		if _, ok := visitedUrls[pageUrl]; !ok {
			pagPrevCount++
			if schema.List.Pagination.Previous.Limit == -1 || pagPrevCount <= schema.List.Pagination.Previous.Limit {
				e.Request.Visit(pageUrl)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit(url)
	return results
}
