package main

// ConnType allows to distinguish between connection type.
type ConnType int

const (
	ConnType_Telnet ConnType = iota
	ConnType_HTTP
)

type UserConn interface {
	HandleBroadcastEvent(BroadcastMessageEvent) error
	GetID() string
	Close()
	GetNickname() string
	IsBlocked(string) bool
	GetConnType() ConnType
}
