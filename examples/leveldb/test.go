package main

import (
	eventstrangler "github.com/balzanelli/event-strangler"
	"log"
)

func main() {
	log.Println("event-strangler/examples/leveldb@v0.1.0")

	config := &eventstrangler.Config{
		InstanceName: "event-strangler/examples/leveldb",
	}

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
