package scrape

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/kuronosu/schema_scraper/pkg/config"
)

type ScrapeDetailsOptions struct {
	Async            bool
	BatchSize        int
	URLs             []string
	Schema           *config.PageSchema
	DetailsCollector DetailsCollector
	Verbose          bool
}

func batchVisit(c *colly.Collector, urls []string, batchSize int) {
	requestCounter := 0
	for _, url := range urls {
		c.Visit(url)
		requestCounter++
		if requestCounter >= batchSize {
			requestCounter = 0
			c.Wait()
		}
	}
	c.Wait()
}

func continuousVisit(c *colly.Collector, urls []string) {
	for _, url := range urls {
		c.Visit(url)
	}
	c.Wait()
}

func ScrapeDetails(options ScrapeDetailsOptions) map[string]string {
	c := colly.NewCollector(colly.Async(options.Async))
	if options.Schema.Cloudflare {
		c.WithTransport(GetCloudFlareRoundTripper())
	}
	mutex := &sync.RWMutex{}

	errorUrls := make(map[string]string)
	responseCounter := 0
	total := len(options.URLs)

	for _, field := range options.Schema.Detail.Fields {
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
				options.DetailsCollector.CollectField(e.Request.URL.String(), field.Name, tmp)
			})
		}(c, field)
	}

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	c.OnScraped(func(r *colly.Response) {
		mutex.Lock()
		responseCounter++
		if options.Verbose {
			fmt.Print("\r[Progress] ", responseCounter, "/", total, " [Errors] ", len(errorUrls))
		}
		mutex.Unlock()
	})

	c.OnError(func(r *colly.Response, err error) {
		mutex.Lock()
		responseCounter++
		errorUrls[r.Request.URL.String()] = err.Error()
		if options.Verbose {
			fmt.Print("\r[Progress] ", responseCounter, "/", total, " [Errors] ", len(errorUrls))
		}
		mutex.Unlock()
	})

	if options.Async {
		batchVisit(c, options.URLs, options.BatchSize)
	} else {
		continuousVisit(c, options.URLs)
	}
	if options.Verbose {
		fmt.Println(" Done")
	}
	return errorUrls
}
