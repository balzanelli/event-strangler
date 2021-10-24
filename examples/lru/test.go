package main

import (
	"encoding/json"
	eventstrangler "github.com/balzanelli/event-strangler"
	"io/ioutil"
	"log"
	"os"
)

func readJson(fileName string) (*map[string]interface{}, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var object map[string]interface{}
	if err = json.Unmarshal(bytes, &object); err != nil {
		return nil, err
	}
	return &object, nil
}

func main() {
	log.Println("event-strangler/examples/lru@v0.1.0")

	store := eventstrangler.NewLRUCacheStore()

	strangler, err := eventstrangler.New(&eventstrangler.Config{
		HashKey: &eventstrangler.HashKeyOptions{
			Name:       "event-strangler/examples/lru",
			Expression: "[subject, transaction_id]",
		},
		Store: store,
	})
	if err != nil {
		log.Fatal(err)
	}

	object, err := readJson("test.json")
	if err != nil {
		log.Fatal(err)
	}

	hashKey, err := strangler.Once(*object)
	if err != nil {
		log.Fatal(err)
	}

	// Run Idempotent Operation

	//if err = strangler.Purge(hashKey); err != nil {
	//	log.Fatal(err)
	//}
	if err = strangler.Complete(hashKey); err != nil {
		log.Fatal(err)
	}
	log.Println(hashKey)
}
