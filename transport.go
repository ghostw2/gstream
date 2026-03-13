package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Consumer interface {
	Start() error
	// Consume () error
}
type Producer interface {
	Start() error
	// Produce()
}

type HTTPConsumer struct {
	listenAddress string
}

func NewHTTPConsumer(listenAddress string) *HTTPConsumer {
	return &HTTPConsumer{listenAddress: listenAddress}
}
func (c *HTTPConsumer) Start() error {
	return http.ListenAndServe(c.listenAddress, c)
}
func (c *HTTPConsumer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}

type HTTPProducer struct {
	listenAddress string
	server        *Server
}

func NewHTTPProducer(s *Server, listenAddress string) *HTTPProducer {
	return &HTTPProducer{
		listenAddress: listenAddress,
		server:        s,
	}
}
func (p *HTTPProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("HTTP transport started ", "port", p.listenAddress)
	path := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println(path)
	paths := strings.Split(path, "/")
	fmt.Println(paths)
	// commit
	if r.Method == http.MethodGet {
		if len(paths) != 2 {
			fmt.Println("Invalid Action")
			return
		}
		topic := paths[1]
		content, err := p.server.topics[topic].Get(0)
		if err != nil {
			fmt.Printf("error getting the message :%v \n ", err)
		}
		fmt.Printf("the content value of the message was :%v \n", string(content))

	}
	if r.Method == http.MethodPost {
		if len(paths) != 2 {
			fmt.Println("Invalid Action")
			return
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		topic := paths[1]
		p.server.createTopic(topic)
		n, err := p.server.topics[topic].Push(content)
		if err != nil {
			fmt.Printf("there was an error producing :%v", err)
		}
		fmt.Printf("successfully pushed :%d", n)
		fmt.Printf("the content pushed was :%v\n", string(content))
	}
}
func (p *HTTPProducer) HandlePublish(w http.ResponseWriter, r *http.Response) {

}
func (p *HTTPProducer) Start() error {
	return http.ListenAndServe(p.listenAddress, p)
}
