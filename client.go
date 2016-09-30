package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

type FindHandler func(string) (Handler, bool)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Client struct {
	send         chan Message
	socket       *websocket.Conn
	findHandler  FindHandler
	dbSession    *r.Session
	stopChannels map[int]chan bool
}

func (c *Client) NewStopChannel(stopkey int) chan bool {
	c.StopForKey(stopkey)

	stop := make(chan bool)
	c.stopChannels[stopkey] = stop
	return stop
}

func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
	}
	delete(c.stopChannels, key)
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			fmt.Println("Err reading JSON:", message)
			fmt.Println("Socket closing: ", client.socket.RemoteAddr())
			break
		}
		if handler, found := client.findHandler(message.Name); found {
			fmt.Println("Calling handler with message:", message)
			handler(client, message.Data)
		}
	}
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		fmt.Println("Read message from send channel:", msg)

		if err := client.socket.WriteJSON(msg); err != nil {
			fmt.Println("Err writing JSON:", msg)
			break
		}
	}
	client.socket.Close()
}

func (c *Client) Close() {
	for _, ch := range c.stopChannels {
		ch <- true
	}
	close(c.send)
}

func NewClient(socket *websocket.Conn, findHandler FindHandler, dbSession *r.Session) *Client {
	return &Client{
		send:         make(chan Message),
		socket:       socket,
		findHandler:  findHandler,
		dbSession:    dbSession,
		stopChannels: make(map[int]chan bool),
	}
}
