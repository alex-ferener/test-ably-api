## Test Ably's API

This (tiny) app will subscribe to a channel, and will publish `Broadcast` messages at regular intervals.
The goal is to calculate the e2e latency between multiple instances (which are using the same channel).

Subscribers will respond with `BroadcastResponse` messages which contains:
- subscriber's instance-id
- the original message
- the difference between the time they received the message and the time the message was published  

You need to set `AUTH_KEY` environment variable.

### Install dependencies and Build

```
make all
```

### Usage

```
./bin/ably-test-amd64-linux --help

Usage:
  -channel string
    	Channel name
  -count uint
    	How many initial messages to send (default 3)
  -duration duration
    	How long to listen for responses (default 30s)
  -instance-id string
    	Instance identifier (must be unique)
  -interval duration
    	How long to wait between sending messages (default 5s)
```

### Example

```
export AUTH_KEY=<REPLACE_ME>

./bin/ably-test-amd64-linux -instance-id instance1 -channel test -count 10 -interval 1s -duration 5s

./bin/ably-test-amd64-linux -instance-id instance2 -channel test -count 10 -interval 1s -duration 5s

InstanceId1     InstanceId2     Min (ms)    Max (ms)    Avg (ms)
-----------     -----------     --------    --------    --------	
instance2       instance3       52          133         87		
instance1       instance3       26          471         187		
instance1       instance2       26          385         213		
```
