/* CodeChat: Collaborative Programming
 * Authors:
 * David Taylor
 * Graham Greving
 */

package CodeChat

import (
	"encoding/json"
	"errors"
	// "flag" // for command line args
	"log"
	"net"
	//	"strings"
)

// Server datatype
type Server struct {
	clients    map[net.Conn]*Client
	numClients int
	// broadcasting channel
	serverChan chan Message
}

// Client datatype
type Client struct {
	server     *Server
	conn       net.Conn
	name       string
	clientChan chan string
}

// IPC Message datatype
type Message struct {
	client Client
	msg    string
	err    error
}

func (serv *Server) broadcast() {
	// loop on incoming messages from the server's chan
	for msg := range serv.serverChan {
		// send message to all clients
		from := serv.clients[msg.client.conn].name + ": "
		i := 0
		for conn, client := range serv.clients {
			if client == &msg.client {
				continue
			}
			to := serv.clients[conn].name
			log.Println("broadcasting to ", to)
			conn.Write([]byte(from))
			conn.Write([]byte(msg.msg))
			i++
		}
		// add support for errors
		log.Println("broadcast ", msg.msg, " to ", i, " clients.")
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

func (serv *Server) newClient(conn net.Conn, name string) (string, error) {
	w := make(chan string)
	c := &Client{serv, conn, name, w}
	serv.clients[conn] = c
	m := name + " entered the chat."
	//serv.serverChan <- msg
	return m, nil
}

func (serv *Server) deleteClient(conn net.Conn) (string, error) {
	name := serv.clients[conn].name
	delete(serv.clients, conn)
	m := name + " " + "has left the chat."
	//serv.serverChan <- msg
	return m, nil
}

func (serv *Server) renameClient(conn net.Conn, newN string, oldN string) (string, error) {
	var m string
	var e error
	if serv.clients[conn].name == oldN {
		// have to delete the current client, and re-add it with the new name
		// kind of lame, and i'm sure there's a better way.
		// if we end up using channels to communicate with the connetions,
		// this will most likely invalidate the channel, so instead we should
		// mutate the name
		serv.deleteClient(conn)
		serv.newClient(conn, newN)
		m = oldN + " " + "now known as" + " " + newN
	} else {
		m = "Failure. Oldname != newname"
		e = errors.New("renameClient: Failed to rename client. oldname != newname")
	}
	return m, e
}

func (client *Client) doCommands(dec *json.Decoder) Message {
	var m Message
	var e error
	var msg string
	m.client = *client
	var v map[string]interface{}
	err := dec.Decode(&v)
	if checkErr(err) {
		goto send
	}
	switch v["cmd"] {
	case "connect":
		if name, ok := v["username"]; ok {
			client.name = name.(string)
		} else {
			e = errors.New("connect: no username given.")
		}
	case "rename":
		if newName, ok := v["newname"]; ok {
			client.name = newName.(string)
		} else {
			e = errors.New("rename: no name(s) given.")
		}
	// case "exit":
	// 	msg, e = serv.deleteClient(conn)
	// 	// expedite the write process so we can kill the connection
	// 	m.msg = msg
	// 	m.err = e
	// 	serv.serverChan <- m
	// 	conn.Close()
	// 	return
	// lots of "msg" this "msg" that. this is a chat message.
	case "msg":
		log.Println("got a mesage")
		if message, ok := v["msg"]; ok {
			msg = message.(string)
		} else {
			e = errors.New("msg: no message given.")
		}
	default:
		e = errors.New("bad JSON given.")
	}
send:
	m.msg = msg
	m.err = e
	return m
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
		m := user.doCommands(dec)
		serv.serverChan <- m
		// user.writeMessage(m)
	}
}

func main() {
	log.Println("CodeChat Server Starting")

	// Initialize the server
	serv := new(Server)
	serv.clients = make(map[net.Conn]*Client)
	serv.serverChan = make(chan Message)

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
