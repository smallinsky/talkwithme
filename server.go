package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/smallinsky/talkwithme/config"
)

func NewServer(cfg *config.Config) (*server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("missing config arg")
	}

	messageStore, err := newFileLogger(cfg.Logfile)
	if err != nil {
		return nil, fmt.Errorf("failed to init message store log: %v", err)
	}

	return &server{
		Clients:      make(map[string]UserConn),
		quitC:        make(chan struct{}),
		eventQueue:   make(chan interface{}),
		config:       cfg,
		messageStore: messageStore,
	}, nil
}

type server struct {
	listener   net.Listener
	Clients    map[string]UserConn
	mtx        sync.RWMutex
	quitC      chan struct{}
	eventQueue chan interface{}

	messageStore messageStore
	config       *config.Config
}

type messageStore interface {
	LogMessage(string) error
	Close() error
}

func (s *server) start() error {
	s.eventLoop()

	if err := s.startTelnetServer(); err != nil {
		return fmt.Errorf("[ERROR] failed to start telnet server: %s", err)
	}
	if err := s.startHTTPServer(); err != nil {
		return fmt.Errorf("[ERROR] failed to start HTTP server: %s", err)
	}
	return nil
}

func (s *server) eventLoop() {
	go func() {
		for {
			select {
			case e := <-s.eventQueue:
				s.handleEvent(e)
			case <-s.quitC:
				return
			}
		}
	}()
}

func (s *server) handleEvent(e interface{}) {
	switch event := e.(type) {
	case BroadcastMessageEvent:
		s.mtx.RLock()
		defer s.mtx.RUnlock()
		for _, client := range s.Clients {
			if client.IsBlocked(event.nickname) {
				continue
			}
			if err := client.HandleBroadcastEvent(event); err != nil {
				log.Printf("[ERROR] handleBroadcastEvent with error: %v", err)
			}
		}
		if err := s.messageStore.LogMessage(event.String()); err != nil {
			log.Printf("[ERROR] failed to store message into store: %v", err)
		}
	case CloseConnEvent:
		s.mtx.Lock()
		defer s.mtx.Unlock()
		delete(s.Clients, event.id)
	default:
		log.Printf("[ERROR] received unknow server event %T\n", e)
	}
}

func (s *server) addUserConn(user UserConn) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.Clients[user.GetID()]; ok {
		return fmt.Errorf("user with ID '%s' already exists", user.GetID())
	}
	s.Clients[user.GetID()] = user
	s.eventQueue <- BroadcastMessageEvent{
		message:       fmt.Sprintf("  ---> %s joined", user.GetNickname()),
		broadcastType: BroadcastType_Info,
		timestamp:     time.Now(),
	}
	return nil
}

func (s *server) close() {
	log.Printf("[INFO] stopping the server\n")
	close(s.quitC)
	if s.listener != nil {
		s.listener.Close()
	}
	s.mtx.Lock()
	for _, c := range s.Clients {
		c.Close()
	}
	s.mtx.Unlock()
	s.messageStore.Close()
}
