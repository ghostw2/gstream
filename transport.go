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
	channel       chan Message
}

func NewHTTPProducer(listenAddress string, channel chan Message) *HTTPProducer {
	return &HTTPProducer{
		listenAddress: listenAddress,
		channel:       channel,
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
		fmt.Println("method get unimplemented")
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
		fmt.Println("messages is being  sent to channel")
		topic := paths[1]
		p.channel <- Message{
			Topic:   topic,
			Content: string(content),
		}
		fmt.Println("messages sent to channel")
	}
}
func (p *HTTPProducer) HandlePublish(w http.ResponseWriter, r *http.Response) {

}
func (p *HTTPProducer) Start() error {
	return http.ListenAndServe(p.listenAddress, p)
}
