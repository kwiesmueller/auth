package ledis

import "github.com/siddontang/goredis"

type client struct {
	client *goredis.Client
}

type Client interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Del(key string) error
	Close()
}

func New(address string, password string) *client {
	c := new(client)
	c.client = goredis.NewClient(address, password)
	return c
}

func (c *client) Ping() error {
	_, err := c.client.Do("PING")
	return err
}

func (c *client) Get(key string) (string, error) {
	return goredis.String(c.client.Do("GET", key))
}

func (c *client) Set(key string, value string) error {
	_, err := c.client.Do("SET", key, value)
	return err
}

func (c *client) Del(key string) error {
	_, err := c.client.Do("DEL", key)
	return err
}

func (c *client) Close() {
	c.client.Close()
}
