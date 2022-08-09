package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/gocolly/colly/v2"
	"github.com/kuronosu/anime_scraper/pkg/config"
	// "gopkg.in/yaml.v2"
	// browser "github.com/EDDYCJY/fake-useragent"
)

const animeflv = "examples/animeflv.yaml"

var urls = []string{
	"https://animeflv.net/anime/the-idolmster",
	"https://animeflv.net/anime/youkoso-jitsuryoku-shijou-shugi-no-kyoushitsu-e-tv-2nd-season",
}

func main() {
	schema, err := config.ReadSchema(animeflv)
	if err != nil {
		log.Fatal(err)
	}
	Scrape(schema, urls)
}

func Scrape(schema *config.PageSchema, urls []string) {
	c := colly.NewCollector()
	if schema.Cloudflare {
		c.WithTransport(GetCloudFlareRoundTripper())
	}
	animes := make(map[string]map[string]interface{})

	for _, field := range schema.Anime.Fields {
		func(c *colly.Collector, field config.Field) {
			c.OnHTML(field.Selector, func(e *colly.HTMLElement) {
				if animes[e.Request.URL.String()] == nil {
					animes[e.Request.URL.String()] = make(map[string]interface{})
				}
				animes[e.Request.URL.String()][field.Name] = field.SafeCompile(e)
			})
		}(c, field)
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	for _, url := range urls {
		c.Visit(url)
	}
	WriteJson(animes, "animes.json")
}

func GetCloudFlareRoundTripper() http.RoundTripper {
	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
	return client.Transport
}

func WriteBody(res *http.Response, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func WriteJson(data interface{}, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	e := json.NewEncoder(outFile)
	e.SetIndent("", "\t")
	return e.Encode(data)
}
