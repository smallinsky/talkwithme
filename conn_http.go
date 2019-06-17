package main

type HTTPClient struct {
	nicname  string
	messageC chan string
	quitC    chan struct{}
	ID       string
}

func (c *HTTPClient) HandleBroadcastEvent(msg BroadcastMessageEvent) error {
	c.messageC <- msg.String()
	return nil
}

func (c *HTTPClient) GetConnType() ConnType {
	return ConnType_HTTP
}

func (c *HTTPClient) Close() {
	close(c.quitC)
}

func (c *HTTPClient) IsBlocked(nickname string) bool {
	return false
}

func (c *HTTPClient) GetNickname() string {
	return c.nicname
}

func (c *HTTPClient) GetID() string {
	return c.ID
}
