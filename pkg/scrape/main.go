package scrape

import (
	"log"

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
				details[e.Request.URL.String()][field.Name] = field.SafeCompile(e)
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

	c.OnHTML(schema.List.ContainerSelector, func(e *colly.HTMLElement) {
		results.Extend(config.ParseLinkData(schema.List.SafeCompile(e)))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit(url)
	return results
}
