package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kuronosu/anime_scraper/pkg/config"
	"github.com/kuronosu/anime_scraper/pkg/scrape"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "",
		Usage: "scrape pages",
		Commands: []*cli.Command{
			{
				Name:    "details",
				Aliases: []string{"d"},
				Usage:   "scrape detail pages",
				Action:  ScrapeDetails,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "schema",
						Value:   "schema.yaml",
						Usage:   "schema that defines the page to be scraped",
						Aliases: []string{"s"},
					},
					&cli.StringFlag{
						Name:    "outfile",
						Value:   "results.json",
						Usage:   "output file",
						Aliases: []string{"o"},
					},
					&cli.StringFlag{
						Name:    "urls",
						Value:   "detail_urls.txt",
						Usage:   "urls to scrape",
						Aliases: []string{"u"},
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "scrape list page",
				Action:  ScrapeList,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "schema",
						Value:   "schema.yaml",
						Usage:   "schema that defines the page to be scraped",
						Aliases: []string{"s"},
					},
					&cli.StringFlag{
						Name:    "outfile",
						Value:   "results.json",
						Usage:   "output file",
						Aliases: []string{"o"},
					},
					&cli.BoolFlag{
						Name:    "details",
						Value:   false,
						Usage:   "scrape detail with the urls in the page",
						Aliases: []string{"d"},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func ScrapeDetails(cCtx *cli.Context) error {
	schema_file := cCtx.String("schema")
	out_file := cCtx.String("outfile")
	urls_file := cCtx.String("urls")

	schema, err := config.ReadSchema(schema_file)
	if err != nil {
		return err
	}
	urls, err := ReadUrls(urls_file)
	if err != nil {
		return err
	}
	details := scrape.ScrapeDetails(schema, urls)
	WriteJson(details, out_file, false)
	return nil
}

func ScrapeList(cCtx *cli.Context) error {
	schema_file := cCtx.String("schema")
	out_file := cCtx.String("outfile")
	scrape_details := cCtx.Bool("details")
	if cCtx.Args().Len() == 0 {
		return fmt.Errorf("no url specified")
	}

	schema, err := config.ReadSchema(schema_file)
	if err != nil {
		return err
	}
	var data interface{}
	list := scrape.ScrapeList(schema, cCtx.Args().Get(0))
	if !scrape_details {
		data = list
	} else {
		_urls := make([]string, len(list))
		i := 0
		for k := range list {
			_urls[i] = k
			i++
		}
		details := scrape.ScrapeDetails(schema, _urls)
		data = details
	}
	WriteJson(data, out_file, false)
	return nil
}

func WriteJson(data interface{}, filename string, indent bool) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	e := json.NewEncoder(outFile)
	if indent {
		e.SetIndent("", "\t")
	}
	return e.Encode(data)
}

func ReadUrls(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	return urls, nil
}
