package main

import (
	"log"

	"github.com/ably/ably-go/ably"
)

func NewApp(config *Config) *App {
	client := newClient(config.AuthKey)
	channel := newChannel(client, config.ChannelName)
	timing := newTiming()

	return &App{
		Config:  config,
		channel: channel,
		timing:  timing,
	}
}

func newClient(authKey string) *ably.Realtime {
	client, err := ably.NewRealtime(ably.WithKey(authKey))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func newChannel(client *ably.Realtime, name string) *ably.RealtimeChannel {
	return client.Channels.Get(name)
}

func main() {
	config := NewConfig()

	app := NewApp(config)
	app.Run()
}
