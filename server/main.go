package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/streadway/amqp"
)

type Post struct {
	Title string
	Url   string
}

type P struct {
	Posts []Post
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Here we connect to RabbitMQ or send a message if there are any errors connecting.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	var posts []Post

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
				// Then we need to add the domain url
				if strings.Contains(postUrl, "item?id=") {
					postUrl = "news.ycombinator.com/" + postUrl
				}
				posts = append(posts, Post{
					Title: postTitle,
					Url:   postUrl,
				})
			}
		})
		file, _ := json.Marshal((P{Posts: posts}))
		_ = ioutil.WriteFile("assets/posts.json", file, 0644)
	})
	c.Visit("https://news.ycombinator.com")
}
