package main

import (
	"fmt"
	"time"
)

const (
	eventTimestampFormat = "15:04:05"
)

type BroadcastType int

const (
	BroadcastType_TextMessage BroadcastType = iota
	BroadcastType_Info
)

type BroadcastMessageEvent struct {
	// source of the event in case of client broadcast event the source is set to client.GetID().
	source string
	// message text payload to broadcast.
	message string
	// nickname of OP
	nickname string
	// room allows to create seperate space from chats.
	room string
	// timestamp when event was created.
	timestamp time.Time
	// broadcastType distingiush beteen clients' events and server notifications.
	broadcastType BroadcastType
}

func (b *BroadcastMessageEvent) String() string {
	receivedAt := b.timestamp.Format(eventTimestampFormat)
	var msg string
	switch b.broadcastType {
	case BroadcastType_TextMessage: // General log format: 20:21:47 <stupefied_adam> message
		msg = fmt.Sprintf("%s <%s> %s", receivedAt, b.nickname, b.message)
	case BroadcastType_Info: // In case of INFO log don't include any username in log entry.
		msg = fmt.Sprintf("%s %s", receivedAt, b.message)
	}
	return msg
}

type CloseConnEvent struct {
	id string
}
