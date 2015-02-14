/* CodeChat: Collaborative Programming
 * Authors:
 * David Taylor
 * Graham Greving
 */

package main

import (
	"fmt"
	"net"
)

// Some sort of a server datatype
type Server struct {
	clients map[net.Conn]Client
	// broadcasting channel
	s_chan chan string
}
// Some sort of a client datatype
type Client struct {
	name   string
	w_chan chan string
}

// This is the message broadcaster; it is run as a goroutine on server
// initialization. It sits and reads the server's chan, which each
// client writes messages to. It then broadcasts the messages back to
// each connection (writes directly to the connection, not message
// passing back to the client)
// Eventually we should ensure that the message is not written back to
// the client that sent the message
func broadcast(serv *Server) {
	for msg := range serv.s_chan {
		for key, _ := range serv.clients {
			// add support for not writing to the client
			// that sent the message
			key.Write([]byte(msg))
		}
		fmt.Print(msg)
	}
}

// Passed an error, logs the error and returns true or false
// Should be used on an if statement to ensure proper termination
// true  -> error
// false -> no error
func checkErr(e error) bool {
	if e != nil {
		fmt.Println("Error", e)
		return true
	}
	return false
}

func log(msg string) {
	fmt.Println(msg)
}

func getClients(serv *Server) {
	for _, value := range serv.clients {
		fmt.Println(value.name)
	}
}

// Connection Handling
func handleConnection(conn net.Conn, serv *Server) {
	// Send a welcome message and read the name
	b := []byte("hey welcome to codechat\n")

	log("new connection!")
	_, err := conn.Write(b)

	// Ensure the connection is closed before this routine exits
	defer conn.Close()

	if checkErr(err) {
		return
	}
	// Read the first line
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if checkErr(err) || n == 0 {
		return
	}
	// Initialize a new Client
	// give it a name and a channel
	// Add it to the server's clients list
	// Eventually, we will replace this with a command dispatcher
	writeChan := make(chan string)
	client := Client{string(buf[0 : n-2]), writeChan}
	serv.clients[conn] = client

	getClients(serv)

	// Now we can handle all of the incoming messages on this client
	for {
		n, err := conn.Read(buf)
		if checkErr(err) || n == 0 {
			break
		}
		//n, err = conn.Write(buf[0:n])
		msg := string(buf[0 : n-2])
		// This is an example of handling commands
		// allows the client to disconnect from the server
		if msg == "exit" {
			delete(serv.clients, conn)
			msg = "Client Left\n"
			getClients(serv)
			serv.s_chan <- msg
			return
		}
		serv.s_chan <- msg + "\n"
		if err != nil {
			fmt.Println(err)
			conn.Close()
			break
		}
	}
}

func main() {
	log("CodeChat Server Starting")

	// Initialize the server
	serv := new(Server)
	serv.clients = make(map[net.Conn]Client)
	serv.s_chan = make(chan string)

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
