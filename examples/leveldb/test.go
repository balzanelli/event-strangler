package main

import (
	eventstrangler "github.com/balzanelli/event-strangler"
	"log"
)

func main() {
	config := &eventstrangler.Config{}

	strangler, err := eventstrangler.NewStrangler(config)
	if err != nil {
		log.Fatal(err)
	}

	_, err = strangler.Once()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("event-strangler@v0.1.0")
}
