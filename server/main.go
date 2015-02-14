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
func server(c <-chan string) {
    for msg := range c {
        fmt.Println(msg)
    }
}

func handleConnection(conn net.Conn, serv *Server, s_chan chan string) {
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
            msg = "Client Left"
            getClients(serv)
            s_chan <- msg
            return
        }
        s_chan <- msg + "\n"
        if err != nil {
            fmt.Println(err)
            conn.Close()
            break
        }
    }
}

type Server struct {
    clients map[net.Conn] Client
}

type Client struct {
    name string
    channel chan string
}

func getClients(serv *Server) {
    for _,value := range serv.clients {
        fmt.Println(value.name)
    }
}

func main() {
    fmt.Println("CodeChat Server Starting")
    s_chan := make(chan string)
    go server(s_chan)
    serv := new(Server)
    serv.clients = make(map[net.Conn] Client)
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
        go handleConnection(conn, serv, s_chan)
    }
}
