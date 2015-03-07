/* CodeChat: Collaborative Programming
 * Authors:
 * David Taylor
 * Graham Greving
 */

package main

import (
	"encoding/json"
	"errors"
	// "flag" // for command line args
	"log"
	"net"
	"sync"
)

// Server datatype
type Server struct {
	clients    map[net.Conn]*Client
	numClients int
	// broadcasting channel
	serverChan chan message
	write      sync.Mutex
}

// Client datatype
type Client struct {
	server     *Server
	conn       net.Conn
	name       string
	clientChan chan string
}

// internal message passing struct for IPC
type message struct {
	client   Client
	msg      OutgoingMessage
	err      error
	res      ClientResponse
	exitflag bool
}

// ClientResponse response message to a client
type ClientResponse struct {
	Success   bool   `json:"success"`
	Cmd       string `json:"cmd"`
	StatusMsg string `json:"status-message"`
}

// OutgoingMessage server message passed to all clients
type OutgoingMessage struct {
	Cmd     string `json:"cmd"`
	From    string `json:"from"`
	Payload string `json:"payload"`
}

func (msg OutgoingMessage) write(conn net.Conn) error {
	log.Println("writing a outgoing message to", conn.RemoteAddr().String())
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Println("OutgoingMessage: marshalled into JSON")
	n, err := conn.Write(b)
	if n == 0 {
		err = errors.New("OutgoingMessage.write: no bytes written")
	}
	log.Println("OutgoingMessage: written")
	return err
}

func (res ClientResponse) write(conn net.Conn) error {
	log.Println("writing a client response to", conn.RemoteAddr().String())
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	n, err := conn.Write(b)
	if n == 0 {
		err = errors.New("ClientResponse.write: no bytes written")
	}
	return err
}

func (serv *Server) broadcast() {
	// loop on incoming messages from the server's chan
	for toBroadcast := range serv.serverChan {
		log.Println("Got a message in broadcast")
		// send message to all clients
		i := 0
		// Have to do connections, not clients. Ask Graham
		for conn := range serv.clients {
			// only write the response to the requesting connection
			if conn == toBroadcast.client.conn {
				toBroadcast.res.write(conn)
			} else {
				// write the message to all other connections
				to := serv.clients[conn].name
				log.Println("broadcasting to ", to)
				toBroadcast.msg.write(conn)
				i++
			}
		}
		log.Println("broadcast to ", i, " clients.")
	}
}

// Passed an error, logs the error and returns true or false
// Should be used on an if statement to ensure proper terserverChanmination
// true  -> error
// false -> no error
func checkErr(e error) bool {
	if e != nil {
		log.Println("checkErr:", e)
		return true
	}
	return false
}

func (serv *Server) getClients(conn net.Conn) (string, error) {
	// Builds an array of names, as well as comma separated string
	// just in case we'll need it later
	names := make([]string, serv.numClients)
	i := 0
	var nameStr string
	for _, value := range serv.clients {
		names[i] = value.name
		nameStr += value.name + ", "
		i++
	}
	return nameStr, nil
}

func (client *Client) doCommands(dec *json.Decoder) (message, error) {
	var m message
	var e error
	var msg string
	var cmd string
	var from string
	m.client = *client
	var v map[string]interface{}
	err := dec.Decode(&v)
	if checkErr(err) {
		return m, err
	}
	from = client.name
	switch v["cmd"] {
	case "connect":
		if name, ok := v["username"]; ok {
			client.name = name.(string)
			msg = client.name
			cmd = "client-connect"
		} else {
			e = errors.New("doCommands: no username passed to connect")
			cmd = "connect"
		}
	case "rename":
		if newName, ok := v["newname"]; ok {
			msg = client.name
			client.name = newName.(string)
			cmd = "client-rename"
			msg += "," + client.name
		} else {
			e = errors.New("doCommands: no name(s) passed to rename")
			cmd = "rename"
		}
	case "exit":
		msg = client.name
		cmd = "client-exit"
		m.exitflag = true
	// lots of "msg" this "msg" that. this is a chat message.
	case "msg":
		log.Println("doCommands: got a mesage")
		if message, ok := v["msg"]; ok {
			msg = from + ": " + message.(string)
			cmd = "message"
		} else {
			e = errors.New("doCommands: no message passed to msg")
			cmd = "message"
		}
	// case "request-write-access":
	// 	client.server.write.Lock()
	// 	cmd = "write-access-granted"
	// case "yield-write-access":
	// 	client.server.write.Unlock()
	// 	cmd = "write-access-yielded"
	case "update-file":
		if file, ok := v["msg"]; ok {
			cmd = "update-file"
			msg = file.(string)
		}
	default:
		e = errors.New("bad JSON given\n in doCommands")
	}
	m.msg = OutgoingMessage{cmd, from, msg}
	// need to fix this errorString
	errorString := ""
	m.res = ClientResponse{e == nil, cmd, errorString}
	m.err = e
	return m, err
}

// Connection Handling
func handleConnection(conn net.Conn, serv *Server) {
	// ensure that the connection is closed before this routine exits
	defer conn.Close()

	log.Println("new connection from " + conn.RemoteAddr().String())

	// Create the client for this connection
	userChan := make(chan string)
	user := &Client{serv, conn, conn.RemoteAddr().String(), userChan}
	serv.clients[conn] = user
	serv.numClients++

	// Create the JSON decoder
	dec := json.NewDecoder(conn)
	for {
		m, err := user.doCommands(dec)
		serv.serverChan <- m
		if m.exitflag || err != nil {
			// write back to clients
			delete(serv.clients, conn)
			serv.numClients--
			return
		}
	}
}

func main() {
	log.Println("CodeChat Server Starting")

	// Initialize the server
	serv := new(Server)
	serv.clients = make(map[net.Conn]*Client)
	serv.serverChan = make(chan message)

	// Start the broadcaster
	go serv.broadcast()

	// Set up networking
	ln, err := net.Listen("tcp", ":8080")
	if checkErr(err) {
		return
	}
	// Handle all incoming connections
	for {
		conn, err := ln.Accept()
		if checkErr(err) {
			break
		}
		go handleConnection(conn, serv)
	}
}
