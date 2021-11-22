package main

import (
	"flag"
	"log"
	"os"
	"time"
)

type Config struct {
	AuthKey         string
	InstanceId      string
	ChannelName     string
	MessageCount    uint
	MessageInterval time.Duration
	TotalDuration   time.Duration
}

func NewConfig() *Config {
	authKey := os.Getenv("AUTH_KEY")
	if authKey == "" {
		log.Fatal("The environment variable 'AUTH_KEY' (API key) needs to be set.")
	}

	instanceId := flag.String("instance-id", "", "Instance identifier (must be unique)")
	channelName := flag.String("channel", "", "Channel name")
	messageCount := flag.Uint("count", 3, "How many initial messages to send")
	interval := flag.Duration("interval", 5*time.Second, "How long to wait between sending messages")
	duration := flag.Duration("duration", 30*time.Second, "How long to listen for responses")

	flag.Parse()

	if *instanceId == "" {
		log.Fatal("-instance-id cannot be empty.")
	}

	if *channelName == "" {
		log.Fatal("-channel cannot be empty.")
	}

	if *messageCount <= 0 {
		log.Fatal("-count must be greater than 0.")
	}

	if *interval <= 0 {
		log.Fatal("-interval must be greater than 0.")
	}

	if *duration <= 0 {
		log.Fatal("-duration must be greater than 0.")
	}

	return &Config{
		AuthKey:         authKey,
		InstanceId:      *instanceId,
		ChannelName:     *channelName,
		MessageCount:    *messageCount,
		MessageInterval: *interval,
		TotalDuration:   *duration,
	}
}
