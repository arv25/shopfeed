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
	stopChannels map[string]chan bool
}

func (c *Client) NewStopChannel(storeId string, channelId string) chan bool {
	stopkey := storeId + "::" + channelId

	// first try to unsub from the channel incase it's already been subscribed to
	c.StopForKey(storeId, channelId)

	// add to stopChannels map
	stop := make(chan bool)
	c.stopChannels[stopkey] = stop
	fmt.Println("Adding stopkey to map: ", stopkey)
	return stop
}

func (c *Client) StopForKey(storeId string, channelId string) {
	stopkey := storeId + "::" + channelId

	if channel, found := c.stopChannels[stopkey]; found {
		fmt.Println("stopping for key: ", stopkey)
		channel <- true
	}

	fmt.Println("Deleting stopkey from map: ", stopkey)
	delete(c.stopChannels, stopkey)
}

func (client *Client) Read() {
	var message Message
	fmt.Println("Listening for data from client on new socket:", client.socket.RemoteAddr())

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
		fmt.Println("Read message from send channel for client to consume:", msg)

		err := client.socket.WriteJSON(msg)
		if err != nil {
			fmt.Println("Err writing JSON:", msg)
			break
		}
		fmt.Println("Putting message on socket")
	}
	client.socket.Close()
}

func (c *Client) Close() {
	for key, ch := range c.stopChannels {
		fmt.Println("Closing routine for: ", key)
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
		stopChannels: make(map[string]chan bool),
	}
}
