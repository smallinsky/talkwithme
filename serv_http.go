package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *server) startHTTPServer() error {
	log.Printf("[INFO] Starting HTTP server on port %s\n", s.config.HTTPPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleMessageGet(w, r)
		case http.MethodPost:
			s.handleMessagePost(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", s.config.HTTPPort),
		Handler: mux,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen and server call failed: %v", err)
		}
	}()

	return nil
}

func (s *server) handleMessageGet(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		// internal error server doesn't support hijacking.
		http.Error(w, "hijacker is not supported", http.StatusInternalServerError)
		return
	}

	c := &HTTPClient{
		nicname:  getRandomName(),
		messageC: make(chan string),
		quitC:    make(chan struct{}),
		ID:       r.RemoteAddr,
	}
	s.addUserConn(c)

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		select {
		// client closed connection
		case <-notify:
			s.eventQueue <- CloseConnEvent{
				id: c.GetID(),
			}
		// server closed connection
		case <-c.quitC:
			return
		}
	}()

	w.Header().Set("Content-Type", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		select {
		case msg := <-c.messageC:
			if _, err := fmt.Fprintln(w, msg); err != nil {
				c.Close()
				return
			}
			fmt.Println("flished http to client ", r.RemoteAddr)
			flusher.Flush()
		case <-c.quitC:
			return
		}
	}
}

type postMessage struct {
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
	// if room is not privided the message will be posted on general.
	Room string `json:"room"`
}

func (m *postMessage) Validate() error {
	if m.Message == "" {
		return fmt.Errorf("message element is required")
	}
	if m.Nickname == "" {
		return fmt.Errorf("nickname is required")
	}
	return nil
}

func (s *server) handleMessagePost(w http.ResponseWriter, r *http.Request) {
	var m postMessage
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid request: proper json payload is required", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := m.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid message: %v", err), http.StatusBadRequest)
		return
	}
	s.eventQueue <- BroadcastMessageEvent{
		message:  m.Message,
		source:   "HTTP_CLIENT",
		nickname: m.Nickname,
		room:     m.Room,
	}
}
