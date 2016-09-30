package main

import (
	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
)

const (
	ChannelMessageStop = iota
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
	stop := client.NewStopChannel(ChannelMessageStop)
	result := make(chan r.ChangeResponse)

	var clientData ChannelSubMsgs
	err := mapstructure.Decode(data, &clientData)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	var compoundIndexQueryVals [2]string
	compoundIndexQueryVals[0] = clientData.StoreId
	compoundIndexQueryVals[1] = clientData.ChannelId

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
	client.StopForKey(ChannelMessageStop)
}

func channelAddMessage(client *Client, data interface{}) {

}
