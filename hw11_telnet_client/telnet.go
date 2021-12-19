package main

import (
	"bufio"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *Client) Close() error {
	err := c.conn.Close()
	return err
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Send() error {
	/* почему не проходит тест test.sh с таким вариантом? при ручном тестировании работает корректно
	buf, _ := io.ReadAll(c.in)
	_, err := c.conn.Write(buf)
	*/

	reader := bufio.NewReader(c.in)
	_, err := reader.WriteTo(c.conn)
	return err
}

func (c *Client) Receive() error {
	buf, _ := io.ReadAll(c.conn)
	_, err := c.out.Write(buf)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
