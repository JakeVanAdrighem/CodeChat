// CodeChat server main file

package main

import (
    "fmt"
    "net"
)

// Some sort of a server datatype

// Some sort of a client datatype

// This is the message broadcaster
// All clients write to this goroutines channel
// To broadcast messages to the server
func server(serv *Server) {
    
    for msg := range serv.s_chan {
        for key,_ := range serv.clients {
            // add support for not writing to the client
            // that sent the message
            key.Write([]byte(msg)) 
        }
        fmt.Println(msg) 
    }
}

func handleConnection(conn net.Conn, serv *Server) {
    b := []byte("hey welcome to codechat\n")
    fmt.Println("new connection!")
    // Read first message from client
    // Parse message for commands
    // Execute commands
    // Write back
    _,err := conn.Write(b)
    if err != nil {
        fmt.Println(err)
        conn.Close()
        return
    }
    buf := make ([]byte, 4096)
    n, err:= conn.Read(buf)
    if err != nil || n == 0 {
        fmt.Println(err)
        conn.Close()
        return
    }
    writeChan := make(chan string)
    client := Client{string(buf[0:n-2]),writeChan}
    serv.clients[conn] = client
    getClients(serv)
    for {
        n, err:= conn.Read(buf)
        if err != nil || n == 0 {
            fmt.Println(err)
            conn.Close()
            break
        }
        //n, err = conn.Write(buf[0:n])
        msg := string(buf[0:n-2])
        if msg == "exit" {
            delete(serv.clients, conn)
            conn.Close()
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

type Server struct {
    clients map[net.Conn] Client
    // broadcasting channel
    s_chan chan string
}

type Client struct {
    name string
    w_chan chan string
}

func getClients(serv *Server) {
    for _,value := range serv.clients {
        fmt.Println(value.name)
    }
}

func main() {
    fmt.Println("CodeChat Server Starting")

    // Initialize the server
    serv := new(Server)
    serv.clients = make(map[net.Conn] Client)
    serv.s_chan = make(chan string)

    // Start the broadcaster
    go server(serv)

    // Set up networking
    ln,err := net.Listen("tcp",":8080")
    if err != nil {
        fmt.Println("Error! Couldn't start server")
        return
    }
    for {
        conn,err := ln.Accept()
        if err != nil {
            fmt.Println("Error! Bad accept.")
            break
        }
        go handleConnection(conn, serv)
    }
}
