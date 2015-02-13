// CodeChat server main file

package main

import (
    "fmt"
    "net"
)

// Some sort of a server datatype

// Some sort of a client datatype

func handleConnection(conn net.Conn) {
    var b =  []byte("hey welcome to codechat\n")
    // Read first message from client
    // Parse message for commands
    // Execute commands
    // Write back
    conn.Write(b)
    conn.Close()
    return
}

func main() {
    fmt.Println("CodeChat Server Starting")

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
        go handleConnection(conn)
    }

}
