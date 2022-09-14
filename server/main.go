package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/streadway/amqp"
)

type Post struct {
	Title string
	Url   string
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

	// We create a Queue to send the message to.
	q, err := ch.QueueDeclare(
		"hn-post-queue", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
	)

	c.OnHTML(".itemlist", func(e *colly.HTMLElement) {
		// We loop over every <tr> item
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

				// Create a JSON object of every post
				b, err := json.Marshal(&Post{Title: postTitle, Url: postUrl})

				// We set the payload for the message.
				body := b
				err = ch.Publish(
					"",     // exchange
					q.Name, // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        body,
					})

				// If there is an error publishing the message, a log will be displayed in the terminal.
				failOnError(err, "Failed to publish a message")
				log.Printf(" [x] sending message: %s", body)
			}
		})
	})
	c.Visit("https://news.ycombinator.com")
}
