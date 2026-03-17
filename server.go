package main

import (
	"fmt"
	"net/http"
)

type Config struct {
	ListenAddr        string
	StoreProducerFunc StoreProducerFunc
}
type Message struct {
	Topic   string
	Content string
}
type Server struct {
	*Config
	topics      map[string]Storer
	consumers   []Consumer
	producers   []Producer
	ProduceChan chan Message
	ConsumeChan chan Message
	quitChannel chan struct{}
}

func NewServer(cfg *Config) (*Server, error) {
	s := &Server{
		Config:      cfg,
		topics:      make(map[string]Storer),
		ProduceChan: make(chan Message, 64),
		ConsumeChan: make(chan Message, 64),
		quitChannel: make(chan struct{}),
	}
	s.producers = []Producer{
		NewHTTPProducer(cfg.ListenAddr, s.ProduceChan),
	}
	return s, nil
}
func (s *Server) loop() {
	for {
		select {
		case msg := <-s.ProduceChan:
			fmt.Println("hello from the loop")
			s.createTopic(msg.Topic)
			_, err := s.topics[msg.Topic].Push([]byte(msg.Content))
			if err != nil {
				fmt.Printf("there was an error pushing to from Producer %v\n", err)
			}
		case <-s.quitChannel:
			fmt.Println("shutting down")
			return
		}
	}
}

func (s *Server) Start() {
	for _, consumer := range s.consumers {
		go func() {
			if err := consumer.Start(); err != nil {
				fmt.Printf("error starting consumer :%v \n", err)
			}
		}()
	}
	for _, p := range s.producers {
		go func() {

			if err := p.Start(); err != nil {
				fmt.Printf("error starting producer :%v \n", err)
			}
		}()
	}
	s.loop()

}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}

func (s *Server) createTopic(name string) bool {
	if _, ok := s.topics[name]; !ok {
		s.topics[name] = s.StoreProducerFunc()
		fmt.Printf("topic created :%v \n", name)
		return true
	}
	return false
}
