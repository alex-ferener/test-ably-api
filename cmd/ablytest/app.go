package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ably/ably-go/ably"
)

const (
	BroadcastEvent         = "Broadcast"
	BroadcastResponseEvent = "BroadcastResponse"
)

type App struct {
	Config  *Config
	channel *ably.RealtimeChannel
	timing  *Timing
}

type BroadcastMsg struct {
	InstanceId string
	Timestamp  int64
}

type BroadcastResponseMsg struct {
	InstanceId  string
	ReceivedMsg BroadcastMsg
	ReceivedAt  int64
}

func (app *App) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), app.Config.TotalDuration)
	defer cancel()

	go app.subscribeBroadcast(ctx)
	go app.subscribeBroadcastResponse(ctx)
	go app.broadcast(ctx)

	// Wait until TotalDuration has passed
	<-ctx.Done()

	// print the timings as a table
	fmt.Print(app.timing)
}

func (app *App) subscribeBroadcast(ctx context.Context) {
	_, err := app.channel.Subscribe(ctx, BroadcastEvent, func(msg *ably.Message) {
		receivedAt := time.Now().UnixNano()
		receivedEvent := BroadcastMsg{}
		err := json.Unmarshal([]byte(msg.Data.(string)), &receivedEvent)
		if err != nil {
			log.Println("Could not decode message", err)
			return
		}

		// Don't respond to your own messages
		// InstanceId is supposed to be unique. If 2 different processes use the same InstanceId the program will misbehave
		if app.Config.InstanceId == receivedEvent.InstanceId {
			return
		}

		err = app.channel.Publish(ctx, BroadcastResponseEvent, BroadcastResponseMsg{
			InstanceId:  app.Config.InstanceId,
			ReceivedMsg: receivedEvent,
			ReceivedAt:  receivedAt,
		})
		if err != nil {
			log.Println("Publish BroadcastResponseEvent error:", err)
		}
	})
	if err != nil {
		log.Println("Broadcast Subscriber error:", err)
	}
}

func (app *App) subscribeBroadcastResponse(ctx context.Context) {
	_, err := app.channel.Subscribe(ctx, BroadcastResponseEvent, func(msg *ably.Message) {
		responseMsg := BroadcastResponseMsg{}
		err := json.Unmarshal([]byte(msg.Data.(string)), &responseMsg)
		if err != nil {
			log.Println("Could not decode message", err)
			return
		}

		log.Printf(
			"A message sent by '%s' at %s was received by '%s' at %s\n",
			responseMsg.ReceivedMsg.InstanceId,
			formatUnix(responseMsg.ReceivedMsg.Timestamp),
			responseMsg.InstanceId,
			formatUnix(responseMsg.ReceivedAt),
		)

		// calculate delta (ns)
		delta := responseMsg.ReceivedAt - responseMsg.ReceivedMsg.Timestamp

		// persist the response
		app.timing.Put(responseMsg.InstanceId, responseMsg.ReceivedMsg.InstanceId, delta)
	})
	if err != nil {
		log.Println("BroadcastResponse Subscriber error:", err)
	}
}

func (app *App) publishBroadcast(ctx context.Context) {
	event := BroadcastMsg{
		InstanceId: app.Config.InstanceId,
		Timestamp:  time.Now().UnixNano(),
	}
	err := app.channel.Publish(ctx, BroadcastEvent, event)

	if err == nil {
		log.Printf("Published BroadcastMsg from %s at %s\n", event.InstanceId, formatUnix(event.Timestamp))
	}
}

func (app *App) broadcast(ctx context.Context) {
	messagesLeft := app.Config.MessageCount

	ticker := time.NewTicker(app.Config.MessageInterval)
	for range ticker.C {
		// stop if no messages left or context deadline exceeded
		if messagesLeft == 0 || ctx.Err() != nil {
			ticker.Stop()
			break
		}
		go app.publishBroadcast(ctx)
		messagesLeft--
	}
}

func formatUnix(unixNano int64) string {
	return fmt.Sprintf("%.6f", float64(unixNano)/float64(time.Second))
}
