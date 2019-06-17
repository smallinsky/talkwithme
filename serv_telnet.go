package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"time"

	"github.com/smallinsky/talkwithme/command"
)

func (s *server) startTelnetServer() error {
	log.Printf("[INFO] Starting telent server on port %s", s.config.TelnetPort)
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.TelnetPort))
	if err != nil {
		return fmt.Errorf("[ERROR] failed to run tcp listener: %v", err)
	}
	s.listener = l

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quitC:
					return
				default:
					fmt.Println("[ERROR] got error during accepting client: ", err)
					continue
				}
			}

			select {
			case <-s.quitC:
				return
			default:
			}

			client := &TelnetClient{
				conn:         conn,
				nickname:     getRandomName(),
				blockedUsers: make(map[string]struct{}),
			}
			s.addUserConn(client)
			go s.handleClient(client)
		}
	}()
	return nil
}

func (s *server) handleClient(client *TelnetClient) {
	r := textproto.NewReader(bufio.NewReader(client.conn))
	for {
		line, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				s.eventQueue <- BroadcastMessageEvent{
					message:       fmt.Sprintf("  <--- %s left", client.GetNickname()),
					broadcastType: BroadcastType_Info,
					timestamp:     time.Now(),
				}
				return
			}
			select {
			case <-s.quitC:
				return
			default:
			}
		}
		cmd, err := command.Scan(line)
		if err != nil {
			if err != command.ErrEmptyLine {
				client.writeErrorMessage(err.Error())
			}
			client.writePrompt()
			continue
		}

		switch t := cmd.(type) {
		case command.TextMessage:
			s.eventQueue <- BroadcastMessageEvent{
				message:   t.Message,
				source:    client.GetID(),
				nickname:  client.nickname,
				timestamp: time.Now(),
			}
		case command.ChangeNickname:
			client.setNickname(t.Nickname)
		case command.BlockUser:
			client.blockUser(t.Nickname)
		case command.UnblockUser:
			client.unblockUser(t.Nickname)
		}
	}
}
