package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 1 * time.Second,
	ReadBufferSize:   100,
	WriteBufferSize:  100,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	messageChan chan []byte
	name        []byte
}

func (client *Client) read() {
	defer func() {
		client.hub.unregister <- client
		fmt.Printf("Client %s disconnected\n", client.name)
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
		if len(client.name) == 0 {
			client.name = message
			client.hub.broadcast <- []byte(fmt.Sprintf("%s joined\n", string(client.name)))
		} else {
			client.hub.broadcast <- bytes.Join([][]byte{client.name, message}, []byte(": "))

		}
	}
}

func (client *Client) write() {
	defer func() {
		fmt.Printf("close message channel to %s\n", client.conn.RemoteAddr().String())
		client.conn.Close()
	}()

	for {
		msg, ok := <-client.messageChan
		if !ok {
			client.conn.WriteMessage(websocket.CloseMessage, []byte("Bye"))
			return
		} else {
			err := client.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Printf("Write error: %v\n", err)
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("new connection from %s\n", conn.RemoteAddr().String())
	client := &Client{
		hub:         hub,
		conn:        conn,
		messageChan: make(chan []byte, 256),
	}
	client.hub.register <- client

	go client.read()
	go client.write()

}
