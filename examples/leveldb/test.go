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

	var dict map[string]interface{}
	if err = json.Unmarshal(bytes, &dict); err != nil {
		return nil, err
	}
	return &dict, nil
}

func main() {
	log.Println("event-strangler/examples/leveldb@v0.1.0")

	config := &eventstrangler.Config{
		InstanceName:      "event-strangler/examples/leveldb",
		HashKeyExpression: "[subject, transaction_id]",
	}

	strangler, err := eventstrangler.NewStrangler(config)
	if err != nil {
		log.Fatal(err)
	}

	dict, err := readJson("test.json")
	if err != nil {
		log.Fatal(err)
	}

	hashKey, err := strangler.Once(*dict)
	if err != nil {
		log.Fatal(err)
	}
	if err = strangler.Complete(hashKey); err != nil {
		log.Fatal(err)
	}
	log.Println(hashKey)
}
