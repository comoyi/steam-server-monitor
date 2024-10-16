package client

import (
	"github.com/comoyi/steam-server-monitor/data"
	"time"
)

func New() *Client {
	client := &Client{}
	return client
}

type Client struct {
	Data *data.Data
}

func (c *Client) Run() {
	go func() {
		for {
			time.Sleep(time.Second)
			c.Data.Counter++
			c.Data.ChCounter <- struct{}{}
		}
	}()
}
