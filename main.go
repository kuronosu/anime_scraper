package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kuronosu/schema_scraper/pkg/config"
	"github.com/kuronosu/schema_scraper/pkg/scrape"
	"github.com/kuronosu/schema_scraper/pkg/utils"
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
					&cli.BoolFlag{
						Name:    "verbose",
						Usage:   "verbose output",
						Value:   true,
						Aliases: []string{"v"},
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
	scrape.SetVerbose(cCtx.Bool("verbose"))

	schema_file := cCtx.String("schema")
	urls_file := cCtx.String("urls")
	outfile := cCtx.String("outfile")
	failed := cCtx.String("failed")

	schema, err := config.ReadSchema(schema_file)
	if err != nil {
		return err
	}
	urls, err := utils.ReadUrls(urls_file)
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
	}
	errors := scrape.ScrapeDetails(options)
	okCount := len(details.Items)
	errCount := len(errors)
	fmt.Println("[OK]", okCount, "[Errors]", errCount, "[Total]", okCount+errCount)
	utils.WriteJson(details.Items, outfile, false)
	if failed != "" {
		utils.WriteJson(errors, failed+".json", false)
		utils.WritePlain(utils.GetErrorUrlsWithoutNotFound(errors), failed+".txt")
	}
	return nil
}

func ScrapeList(cCtx *cli.Context) error {
	scrape.SetVerbose(cCtx.Bool("verbose"))

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
		return utils.WritePlain(list, outfile+".txt")
	}
	list := scrape.ScrapeList(schema, cCtx.Args().Get(0))
	return utils.WriteJson(list, outfile+".json", false)
}
