package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/centrifugal/centrifuge-go"
)

type WSClient struct {
	client *centrifuge.Client
}

func NewWSClient(client *centrifuge.Client) *WSClient {
	l := &WSClient{client: client}
	client.OnError(l)
	return l
}

func (l *WSClient) Run(channel string) error {
	sub, err := l.client.NewSubscription(channel)
	if err == nil {
		sub.OnPublish(l)
		sub.OnSubscribeError(l)
	}
	return err
}

func (l *WSClient) OnSubscribeError(sub *centrifuge.Subscription, event centrifuge.SubscribeErrorEvent) {
	panic(fmt.Errorf(event.Error))
}

func (l *WSClient) OnPublish(sub *centrifuge.Subscription, event centrifuge.PublishEvent) {
	fmt.Println(string(event.Data))
}

func (l *WSClient) OnError(client *centrifuge.Client, event centrifuge.ErrorEvent) {
	panic(fmt.Errorf(event.Message))
}

func main() {
	var u = flag.String(
		"u",
		"ws://127.0.0.1:8000/connection/websocket",
		"url for centrifugo",
	)

	var c = flag.String(
		"c",
		"channel",
		"channel name",
	)

	fmt.Println("url:", *u)
	fmt.Println("channel:", *c)

	centconf := centrifuge.DefaultConfig()
	centconf.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	cf := centrifuge.New(*u, centconf)
	if err := NewWSClient(cf).Run(*c); err != nil {
		panic(err)
	}

	if err := cf.Connect(); err != nil {
		panic(err)
	}
	defer cf.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt
}
