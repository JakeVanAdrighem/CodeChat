// Simple client for CodeChat server

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Message struct {
	Cmd string `json:"cmd"`
	Msg string `json:"msg"`
}

type Connect struct {
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
}

func read(conn net.Conn) {
	// b := make([]byte, 4096)
	d := json.NewDecoder(conn)
	for {
		var v map[string]interface{}
		err := d.Decode(&v)
		if err != nil {
			log.Println("error, bad json")
			return
		}
		// Catch a response from the server
		if s, ok := v["success"]; ok {
			if s.(bool) {
				log.Println("success")
				continue
			} else {
				log.Println("read: command failed")
				log.Println("returned: ", v["status-message"])
			}
			continue
		} else if c, ok := v["cmd"]; ok { // Catch general messages from the server
			switch c {
			case "message":
				log.Println("got message")
			case "client-exit":
				log.Println("client exited")
			case "client-connect":
				log.Println("client entered")
			default:
				log.Println("no cmd parsed")
			}
		}

		// n, e := conn.Read(b)
		// if e != nil || n == 0 {
		// 	log.Println("error.")
		// 	return
		// }
		// fmt.Println(string(b))
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("need a name")
		return
	}
	log.Println(args[0])

	c, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	go read(c)

	user := Connect{"connect", args[0]}
	b, err := json.Marshal(user)

	n, err := c.Write(b)
	if err != nil || n == 0 {
		log.Println(err)
		return
	}
	// keep the connection alive
	for {
		read := make([]byte, 4096)
		n, err := os.Stdin.Read(read)
		if err != nil || n == 0 {
			log.Println(err)
			return
		}
		readStr := strings.TrimSpace(string(read))
		m := Message{"msg", readStr}
		//fmt.Println(m)
		b, e := json.Marshal(m)
		if e != nil {
			log.Println("somethin happened...")
			continue
		}
		c.Write(b)
		//fmt.Println(b)
	}
}
