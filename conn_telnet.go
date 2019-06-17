package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/smallinsky/talkwithme/vt"
)

type TelnetClient struct {
	ID string
	// nickname of ther user, when user connects to the server random one is generated
	// but /nicname client command allows change this property.
	nickname string

	// list of bloced users. Right now users are banned per nickname,
	// this is not save solution because banned user can easily change their nickname.
	// to bypass block mechanism, the best solution  is to introduce user nickname registration
	// functionality and implement authentication mechanism into server logic.
	blockedUsers map[string]struct{}

	conn net.Conn
	mtx  sync.RWMutex
}

func (c *TelnetClient) HandleBroadcastEvent(event BroadcastMessageEvent) error {
	if event.source == c.GetID() {
		event.nickname = vt.Cyan(event.nickname)
		// delete raw line from client terminal to make a space for full formatted log entry.
		c.conn.Write([]byte(vt.MoveCurUp(1) + vt.DeleteNLines(1)))
	}

	text := event.String()
	if event.broadcastType == BroadcastType_Info {
		text = vt.Yellow(text)
	}

	msg := fmt.Sprintf(vt.ClearLineCurLeft()+"%s\n%s> ", text, c.GetNickname())
	_, err := c.conn.Write([]byte(msg))
	return err
}

func (c *TelnetClient) writeErrorMessage(errMsg string) error {
	msg := fmt.Sprintf("%s\n", vt.Red(errMsg))
	if _, err := c.conn.Write([]byte(msg)); err != nil {
		return err

	}
	return nil

}

func (c *TelnetClient) IsBlocked(nickname string) bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	_, blocked := c.blockedUsers[nickname]
	return blocked
}

func (c *TelnetClient) blockUser(nickname string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.blockedUsers[nickname] = struct{}{}
}

func (c *TelnetClient) unblockUser(nickname string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	delete(c.blockedUsers, nickname)
}

func (c *TelnetClient) setNickname(newNicname string) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	c.nickname = newNicname
}

func (c *TelnetClient) GetNickname() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.nickname
}

func (c *TelnetClient) GetID() string {
	return c.conn.RemoteAddr().String()
}

func (c *TelnetClient) GetConnType() ConnType {
	return ConnType_Telnet
}

func (c *TelnetClient) writePrompt() error {
	_, err := c.conn.Write([]byte(fmt.Sprintf("%s> ", c.nickname)))
	if err != nil {
		return err
	}
	return nil
}

func (c *TelnetClient) Close() {
	c.conn.Close()
}
