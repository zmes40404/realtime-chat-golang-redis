package chat

/*
 Encapsulation of each user connection (including sending and receiving channels)
1. Association with Hub
2. WebSocket entity
3. Send/Receive Messages
4. Goroutine collaboration mechanism
*/

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	_"golang.org/x/text/message"
)

// Avoid infinite WebSocket waits and ensure resources are released
// Server ping-> Client
// Client pong-> Server
const (

	// The timeout period for WebSocket to wait when writing messages (server writing back to client).
	// Why setting? -> 1. Prevent goroutine from being stuck due to a client network jam. 2. If the write cannot be completed within 10 seconds, the connection will be automatically closed.
	// 5-30 seconds is usually recommended
	writeWait = 10 * time.Second	// timeout of writing

	// Maximum wait time for a pong response (for keepalive)
	// If no pong is received within this period, the client is considered disconnected and the connection is automatically closed.
	// Commonly setting 30~120 seconds
	pongWait = 60 * time.Second

	// "pingPeriod": How often to send a ping to client, if client receive the ping it will return a pong (shorter than pongWait), pingPeriod < pongWait: The server will periodically check if the client is still alive
	// General setting: 1. should smaller than pongWait 2.80%~90% of pongWait is the best
	// If set too close to pongWait, a race condition may occur: the server timeout is detected before the pong arrives.
	pingPeriod = (pongWait * 9) / 10	

	// maxMessageSize setting: (Why setting? -> 1. Prevent clients from maliciously sending oversized messages and causing memory overflow (DoS) 2. Limiting message size can also simplify back-end data analysis and processing.)
	// 1. Chat room: Generally, a limit of 256~1024 bytes is sufficient. 2. Multi-person collaborative notes: It may be set to a larger size, such as 2KB, 4KB
	maxMessageSize = 512	// Maximum message size -> Prevent chat messages from crashing the server when someone uploads 100MB(Limit the maximum byte size of messages sent from the client)
)

// Client represents a single chatting user
type Client struct {
	Conn *websocket.Conn	// For websocket connection
	Send chan []byte	// Message to be sent (sent by the server), A channel used to send messages to the client (written by the hub and read by the client goroutine)
	Room *Room		// Chat room
}


// "ReadPump" read messages from WebSocket connection, and responsible for receiving messages sent by users
func (c *Client) ReadPump() {
	defer func() {
		c.Room.removeClient(c)		// Remove myself from the Room
		c.Conn.Close()		// Close connection
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	// keep-alive logic
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))	// If no message (including Pong) is received for a pongWait, the connection will be closed by timeout
	c.Conn.SetPongHandler(func(string) error{		// When a Pong is received, reset the deadline
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {	// Always read messages (WebSocket Text / Binary)
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break	// Once the read fails (for example, the network is disconnected or the connection is closed), the resource recovery process will be terminated.
		}

		// Broadcast to all clients in the same room
		c.Room.broadcast(message, c)
	}
}

// "WritePump" writes messages to the WebSocket connection(to users)
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)	// Set a Ticker to send a Ping every once in a while
	defer func() {		// Cleanup logic: stop ticker + close connection
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:		// Receive a message from the Send Channel to send out, 訊息先透過 channel 傳遞，避免race condition
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				// The channel is closed by the room, indicating that the connection is about to be disconnected
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} 

			w, err := c.Conn.NextWriter(websocket.TextMessage)		// Use NextWriter() to start a TextMessage writer
			if err != nil {
				return
			}
			w.Write(message)

			// Combine queued messages, write out all the information accumulated during this period
			n := len(c.Send)	// Process the messages in the buffer that have not been sent out (if too many people are talking at the same time)
			for i:=0; i<n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return 
			}


		case <-c.Send:
		// case message, ok := <-c.Send:
			// Send pings regularly to maintain active connections
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return 
			}
		}
	}
}