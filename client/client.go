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
		if s, ok := v["success"]; ok && !s.(bool) {
			log.Println("read: command failed")
			log.Println("returned: ", v["status-message"])
			continue
		} else {
		    from := strings.TrimSpace(v["from"].(string))
			msg := strings.TrimSpace(v["status-message"].(string))
			fmt.Println(from + ": " + msg)
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
