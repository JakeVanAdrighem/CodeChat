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
	clients map[net.Conn]Client
	// broadcasting channel
	serverChan chan string
}

// Client datatype
type Client struct {
	conn net.Conn
	name string
	clientChan chan string
}

func broadcast(serv *Server) {
	// loop on incoming messages from the servers chan
	for msg := range serv.serverChan {
		// send message to all clients
		for key := range serv.clients {
			// add support for not writing to the client
			// that sent the message
			key.Write([]byte(msg))
		}
		log.Println(msg)
	}
}

// Passed an error, log.Printlns the error and returns true or false
// Should be used on an if statement to ensure proper termination
// true  -> error
// false -> no error
func checkErr(e error) bool {
	if e != nil {
		log.Println("Error", e)
		return true
	}
	return false
}

func getClients(serv *Server) {
	for _, value := range serv.clients {
		log.Println(value.name)
	}
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
		for key,value := range(v) {
			log.Println(key + ":", value)
		}
		if v["cmd"] != nil {
			log.Println("Got a cmd")
		}
	}
}

func main() {
	log.Println("CodeChat Server Starting")

	// Initialize the server
	serv := new(Server)
	serv.clients = make(map[net.Conn]Client)
	serv.serverChan = make(chan string)

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
