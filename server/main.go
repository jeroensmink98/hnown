package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Post struct {
	Title string
	Url   string
}

type P struct {
	Posts []Post
}

func main() {
	var posts []Post

	fName := "assets/data.csv"
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

				// Check if post is either a show HN or Ask HN type of post
				if strings.Contains(postUrl, "item?id=") {

				}

				posts = append(posts, Post{
					Title: postTitle,
					Url:   postUrl,
				})

			}
		})
		myJson, _ := json.Marshal((P{Posts: posts}))
		fmt.Println(string(myJson))
	})

	c.Visit("https://news.ycombinator.com")
}
