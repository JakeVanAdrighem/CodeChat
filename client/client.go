// client.go
// Authors: Graham Greving, David Taylor, Jake VanAdrighem
// CMPS 112: Final Project - CodeChat

// This program slowly evolved from a simple way to test the server
// to a fully fledged client program, and then into a library/package
// which is exportable to any other go program. It exports the structs:
// 		ReadMessage
//		WriteMessage
//		Client
//  (it technically exports the other ones, but only so that they can
//   be marshalled into json)
// It also exports the following functions:
// 		Client.Read()
// 		Client.Write()
// 		Connect

package client

import (
	"encoding/json"
	//"flag"
	//"fmt"
	"errors"
	//"log"
	"net"
	//"os"
	//"strings"
)

// WriteMessage : A message written to the server. Contains a command
// and a message. The language of this struct intends for it to be an
// actual chat message, but it can be overloaded to send arbitrary
// commands to the server. This is an example of it's evolution from
// a simple client to a library.
type WriteMessage struct {
	Cmd string `json:"cmd"`
	Msg string `json:"msg"`
}
// ReadMessage : A message read from the server. Contains a command
// who it's from, and a payload.
type ReadMessage struct {
	Cmd string
	From string
	Payload string
}

// ConnectMsg
// Doesn't need to be visible to any program importing this package,
// only visible because of json marshalling
type ConnectMsg struct {
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
}

// Client : The Client datatype encapsulates any connections and data
// associated with a connection to the CodeChat server: username, the
// netowork connection, and the JSON decoder. Also included is a write
// access flag which is in testing.
type Client struct {
	Username string
	Conn	net.Conn
	Jreader *json.Decoder
	WriteAccess	bool
}

// Read reads in a message from the server, catches any errors and
// repackages it as a nice exported struct for easy handling server
// commands
// It is recommended that this function is called inside an infinite 
// for-loop inside a goroutine.
func (client *Client) Read() (ReadMessage, error){
	var v map[string]interface{}
	var retMsg ReadMessage
	err := client.Jreader.Decode(&v)
	if err != nil {
		// should exit here : signifies a dead server
		return retMsg, err
	}
	// Catch a response from the server
	if s, ok := v["success"]; ok {
		if s.(bool) {
			retMsg.Cmd = "success"
			retMsg.From = "server"
		} else {
			err = errors.New("Read: previous command failed")
		}
		// Catch general messages from the server
	} else if c, ok := v["cmd"]; ok {
		switch c {
		case "message":
			retMsg.Cmd = c.(string)
		case "client-exit":
			retMsg.Cmd = c.(string)
		case "client-connect":
			retMsg.Cmd = c.(string)
		case "update-file":
			retMsg.Cmd = c.(string)
		case "request-write-access":
		case "yield-write-access":
		default:
			// catches when a client EOFs
			err = errors.New("Read: no cmd parsed")
		}
		if f, ok := v["from"]; ok {
			retMsg.From = f.(string)
		} else {
			err = errors.New("Read: no from given")
		}
		if p, ok := v["payload"]; ok {
			// acceptable error
			retMsg.Payload = p.(string)
		} else {
			// acceptable error
			err = errors.New("Read: no payload given")
		}
	} else {
		err = errors.New("Read: json parsing failed")
	}
	return retMsg, err
}

// Write : Writes a command and message to the sever.
func (c *Client) Write(command string, payload string) error {
	m := WriteMessage{command, payload}
	var err error
	b, e := json.Marshal(m)
	if e != nil {
		err = errors.New("Write: json marshal failed: " + e.Error())
		return e
	}
	n, e := c.Conn.Write(b)
	if n == 0 || e != nil {
		err = errors.New("Write: conn.Write failed: " + e.Error())
	}
	return err
}

// Close : Write an exit command and close off the connection
// Should be called with a defer to ensure it happens.
// defer client.Close("reason")
func (c *Client) Close(reason string) {
	c.Write("exit", reason)
	c.Conn.Close()
}

// Connect : Connects a user to the server, does all the necessary
// networking and connecting with the given username. Also initializes
// the JSON decoder and the writeaccess flag.
// Returns the Client. It is up to the user to close the client's
// connection. Simply calling defer Client.Conn.Close() will do the
// trick.
func Connect(username string) (*Client, error) {
	var c = new(Client)
	var err error
	c.Username = username
	// make the connection to the server
	c.Conn, err = net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}
	c.Jreader = json.NewDecoder(c.Conn)
	user := ConnectMsg{"connect", c.Username}
	
	b, err := json.Marshal(user)
	n, err := c.Conn.Write(b)
	if err != nil || n == 0 {
		return nil, err
	}
	c.WriteAccess = false
	return c, err
}

