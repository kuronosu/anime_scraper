package scrape

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/kuronosu/anime_scraper/pkg/config"
)

func ScrapeDetails(schema *config.PageSchema, urls []string) map[string]map[string]interface{} {
	c := colly.NewCollector()
	if schema.Cloudflare {
		c.WithTransport(GetCloudFlareRoundTripper())
	}
	details := make(map[string]map[string]interface{})

	for _, field := range schema.Detail.Fields {
		func(c *colly.Collector, field config.Field) {
			c.OnHTML(field.Selector, func(e *colly.HTMLElement) {
				if details[e.Request.URL.String()] == nil {
					details[e.Request.URL.String()] = make(map[string]interface{})
				}
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
				details[e.Request.URL.String()][field.Name] = tmp
			})
		}(c, field)
	}

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	for _, url := range urls {
		c.Visit(url)
	}
	return details
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
