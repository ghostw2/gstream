package main

import (
	"fmt"
	"net/http"
)

type Config struct {
	ListenAddr        string
	StoreProducerFunc StoreProducerFunc
}

type Server struct {
	*Config
	topics      map[string]Storer
	consumers   []Consumer
	producers   []Producer
	quitChannel chan struct{}
}

func NewServer(cfg *Config) (*Server, error) {
	s := &Server{
		Config:      cfg,
		topics:      make(map[string]Storer),
		quitChannel: make(chan struct{}),
	}
	s.producers = []Producer{
		NewHTTPProducer(s, cfg.ListenAddr),
	}
	return s, nil
}

func (s *Server) Start() {
	for _, consumer := range s.consumers {
		if err := consumer.Start(); err != nil {
			fmt.Printf("error starting consumer :%v \n", err)
		}
	}
	for _, p := range s.producers {
		if err := p.Start(); err != nil {
			fmt.Printf("error starting producer :%v \n", err)
		}
	}
	<-s.quitChannel
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}

func (s *Server) createTopic(name string) bool {
	if _, ok := s.topics[name]; !ok {
		s.topics[name] = s.StoreProducerFunc()
		return true
	}
	return false
}
