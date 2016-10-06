package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
)

func channelList(client *Client, data interface{}) {
	cursor, err := r.Table("channels").
		Changes(r.ChangesOpts{IncludeInitial: true}).
		Run(client.dbSession)
	if err != nil {
		client.send <- Message{"error", err.Error()}
	}

	go func() {
		var change r.ChangeResponse
		for cursor.Next(&change) {
			client.send <- Message{"channel", change.NewValue}
		}
		cursor.Close()
	}()
}

func channelSubscribeMessages(client *Client, data interface{}) {
	result := make(chan r.ChangeResponse)

	var clientData ChannelSubMsgs
	err := mapstructure.Decode(data, &clientData)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	// create a new stop channel for this feed
	stop := client.NewStopChannel(clientData.StoreId, clientData.ChannelId)

	var compoundIndexQueryVals [2]string
	compoundIndexQueryVals[0] = clientData.StoreId
	compoundIndexQueryVals[1] = clientData.ChannelId

	fmt.Println("Cursor opening for channel:", clientData.ChannelId)
	cursor, err := r.Table("messages").
		GetAll(compoundIndexQueryVals).
		OptArgs(r.GetAllOpts{Index: "StoreChannel"}).
		Changes(r.ChangesOpts{IncludeInitial: true}).
		Run(client.dbSession)
	if err != nil {
		client.send <- Message{"error", err.Error()}
	}

	go func() {
		var change r.ChangeResponse
		for cursor.Next(&change) {
			result <- change
		}
	}()

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Closing DB cursor for channel:", clientData.ChannelId)
				cursor.Close()
				return
			case change := <-result:
				if change.NewValue != nil && change.OldValue == nil {
					client.send <- Message{"channel message", change.NewValue}
				}
			}
		}
	}()
}

func channelUnsubscribeMessages(client *Client, data interface{}) {
	var clientData ChannelSubMsgs
	err := mapstructure.Decode(data, &clientData)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	// Pass the store_id and channel_id to determine which key in the map of stop channels.
	client.StopForKey(clientData.StoreId, clientData.ChannelId)
}

func channelAddMessage(client *Client, data interface{}) {
	var clientData ChannelAddMsg

	if err := mapstructure.Decode(data, &clientData); err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	if err := r.Table("messages").Insert(clientData).Exec(client.dbSession); err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	client.send <- Message{"Got message from client-side", clientData}
}
