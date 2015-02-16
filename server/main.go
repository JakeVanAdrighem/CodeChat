/* CodeChat: Collaborative Programming
 * Authors:
 * David Taylor
 * Graham Greving
 */

package main

import (
	"encoding/json"
	"log"
	"net"
)

// Server datatype
type Server struct {
	clients    map[net.Conn]Client
	numClients int
	// broadcasting channel
	serverChan chan Message
}

// Client datatype
type Client struct {
	conn       net.Conn
	name       string
	clientChan chan string
}

// IPC Message datatype
type Message struct {
	conn net.Conn
	msg  string
}

func broadcast(serv *Server) {
	// loop on incoming messages from the servers chan
	for msg := range serv.serverChan {
		// send message to all clients
		for client := range serv.clients {
			// add support for not writing to the client
			// that sent the message
			client.Write([]byte(msg.msg))
		}
		log.Println(msg)
	}
}

// Passed an error, logs the error and returns true or false
// Should be used on an if statement to ensure proper termination
// true  -> error
// false -> no error
func checkErr(e error) bool {
	if e != nil {
		log.Println(e)
		return true
	}
	return false
}

func getClients(serv *Server, conn net.Conn) {
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
	msg := Message{conn, nameStr}
	serv.serverChan <- msg
	log.Println(names)
}

func newClient(serv *Server, conn net.Conn, name string) {
	w := make(chan string)
	c := Client{conn, name, w}
	serv.clients[conn] = c
	msg := Message{conn, name + " " + "entered the chat."}
	serv.serverChan <- msg
}

func deleteClient(serv *Server, conn net.Conn) {
	name := serv.clients[conn].name
	delete(serv.clients, conn)
	msg := Message{conn, name + " " + "has left the chat."}
	serv.serverChan <- msg
}

// Connection Handling
func handleConnection(conn net.Conn, serv *Server) {
	// Send a welcome message and read the name
	b := []byte("hey welcome to codechat\n")

	_, err := conn.Write(b)
	if checkErr(err) {
		return
	}
	log.Println("new connection!")
	// Ensure the connection is closed before this routine exits
	defer conn.Close()
	dec := json.NewDecoder(conn)
	// Now we can handle all of the incoming messages on this client
	for {
		var v map[string]interface{}
		err := dec.Decode(&v)
		if checkErr(err) {
			break
		}
		// log all of the json args:
		for key, value := range v {
			log.Println(key+":", value)
		}
		if v["cmd"] != nil {
			log.Println("Got a cmd: ", v["cmd"])
		}
		switch v["cmd"] {
		case "connect":
			if v["username"] != nil {
				log.Println("new user:", v["username"].(string))
				newClient(serv, conn, v["username"].(string))
			} else {
				log.Println("no username given for connect cmd.")
			}
		case "rename":
			if v["oldname"] != nil || v["newname"] != nil {

			}
		case "exit":
			deleteClient(serv, conn)
		default:
			log.Println("Bad JSON recieved")
		}
	}
}

func main() {
	log.Println("CodeChat Server Starting")

	// Initialize the server
	serv := new(Server)
	serv.clients = make(map[net.Conn]Client)
	serv.serverChan = make(chan Message)

	// Start the broadcaster
	go broadcast(serv)

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
