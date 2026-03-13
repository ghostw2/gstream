package main

import (
	"log"
)

func main() {
	cfg := &Config{
		ListenAddr: ":3000",
		StoreProducerFunc: func() Storer {
			return NewMemoryStore()
		},
	}
	s, err := NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// s.producers = []Producer{
	// 	NewHTTPProducer(":3000"),
	// 	NewHTTPConsumer(":4000"),
	// }
	s.Start()

}
