// Simple client for CodeChat server

package main

import (
	"log"
	"net"
)

func main() {
	c,err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	b := []byte(`{"cmd":"connect","username":"gramasaurous"}`)
	log.Println(string(b))
	n,err := c.Write(b)
	if err != nil || n == 0 {
		log.Println(err)
		return
	}
	// keep the connection alive
	for {}
}
