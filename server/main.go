package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type post struct {
	title string
	url   string
}

func newPost(title string, url string) *post {
	p := post{}
	p.title = title
	p.url = url
	return &p
}

func main() {
	fName := "data.csv"
	file, err := os.Create(fName)

	if err != nil {
		log.Fatalf("Could not create file, err: %q", err)
		return
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
	)

	c.OnHTML(".itemlist", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {

			// Check if the td item is not emtpy
			if len(el.ChildText("td:nth-child(3)")) != 0 {
				postTitle := el.ChildText("td:nth-child(3) > a")
				postUrl := el.ChildAttr("td:nth-child(3) > a", "href")

				writer.Write([]string{
					postTitle,
					postUrl,
				})
			}
		})
	})

	c.Visit("https://news.ycombinator.com")
}
