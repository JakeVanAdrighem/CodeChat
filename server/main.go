// CodeChat server main file

package main

import (
    "fmt"
    "net"
)

// Some sort of a server datatype

// Some sort of a client datatype

func server(c <-chan string) {
    for msg := range c {
        fmt.Println(msg)
    }
}

func handleConnection(conn net.Conn, s_chan chan string) {
    var b =  []byte("hey welcome to codechat\n")
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
            conn.Close()
            msg = "Client Left"
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

func main() {
    fmt.Println("CodeChat Server Starting")
    s_chan := make(chan string)
    go server(s_chan)
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
    
        go handleConnection(conn, s_chan)
    }
}
