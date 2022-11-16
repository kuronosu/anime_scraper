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
						Value:   "urls.txt",
						Usage:   "urls to scrape",
						Aliases: []string{"u"},
					},
					&cli.BoolFlag{
						Name:  "async",
						Value: false,
						Usage: "scrape pages asynchronously",
					},
					&cli.StringFlag{
						Name:    "failed",
						Usage:   "output file for failed urls without extension, if it is not empty",
						Value:   "",
						Aliases: []string{"f"},
					},
					&cli.IntFlag{
						Name:    "batch-size",
						Usage:   "batch size for scraping, affects only when async is true",
						Value:   500,
						Aliases: []string{"b"},
					},
					&cli.BoolFlag{
						Name:    "verbose",
						Usage:   "verbose output",
						Value:   true,
						Aliases: []string{"v"},
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
						Value:   "urls",
						Usage:   "output file, without extension",
						Aliases: []string{"o"},
					},
					&cli.BoolFlag{
						Name:    "flat",
						Value:   true,
						Usage:   "flatten the list",
						Aliases: []string{"f"},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getErrorUrlsWithoutNotFound(errors map[string]string) []string {
	var urls []string
	for url, err := range errors {
		if err != "Not Found" {
			urls = append(urls, url)
		}
	}
	return urls
}

func ScrapeDetails(cCtx *cli.Context) error {
	schema_file := cCtx.String("schema")
	urls_file := cCtx.String("urls")
	outfile := cCtx.String("outfile")
	failed := cCtx.String("failed")

	schema, err := config.ReadSchema(schema_file)
	if err != nil {
		return err
	}
	urls, err := ReadUrls(urls_file)
	if err != nil {
		return err
	}
	details := scrape.NewMemoryDetailsCollector()
	options := scrape.ScrapeDetailsOptions{
		Async:            cCtx.Bool("async"),
		Schema:           schema,
		URLs:             urls,
		DetailsCollector: details,
		BatchSize:        cCtx.Int("batch-size"),
		Verbose:          cCtx.Bool("verbose"),
	}
	errors := scrape.ScrapeDetails(options)
	okCount := len(details.Items)
	errCount := len(errors)
	fmt.Println("[OK]", okCount, "[Errors]", errCount, "[Total]", okCount+errCount)
	WriteJson(details.Items, outfile, false)
	if failed != "" {
		WriteJson(errors, failed+".json", false)
		WritePlain(getErrorUrlsWithoutNotFound(errors), failed+".txt")
	}
	return nil
}

func ScrapeList(cCtx *cli.Context) error {
	schema_file := cCtx.String("schema")
	outfile := cCtx.String("outfile")
	flat := cCtx.Bool("flat")
	if cCtx.Args().Len() == 0 {
		return fmt.Errorf("no url specified")
	}

	schema, err := config.ReadSchema(schema_file)
	if err != nil {
		return err
	}
	if flat {
		list := scrape.ScrapeListFlat(schema, cCtx.Args().Get(0))
		return WritePlain(list, outfile+".txt")
	}
	list := scrape.ScrapeList(schema, cCtx.Args().Get(0))
	return WriteJson(list, outfile+".json", false)
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

func WritePlain(data []string, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	for _, str := range data {
		outFile.WriteString(str + "\n")
	}
	return nil
}
