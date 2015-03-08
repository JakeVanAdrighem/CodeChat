// Simple client for CodeChat server

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

type WriteMessage struct {
	Cmd string `json:"cmd"`
	Msg string `json:"msg"`
}

type ReadMessage struct {
	Cmd string
	From string
	Payload string
}

type ConnectMsg struct {
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
}

type Client struct {
	Username string
	Conn	net.Conn
	Jreader *json.Decoder
	WriteAccess	bool
}


// Read reads in a message from the server, catches any errors and
// repackages it as a nice exported struct for easy handling server
// commands
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

