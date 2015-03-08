package main

import (
	codechat "CodeChat/client"
	"log"
	"os"
	"strings"
	//"encoding/json"
)

func runReader(client *codechat.Client) {
	var read codechat.ReadMessage 
	for {
		var err error
		read, err = client.Read()
		if err != nil {
			log.Println(err)
			if err.Error() == "EOF" {
				return
			}
		}
		switch read.Cmd {
		case "success":
			log.Println("success")
		default:
			log.Println(read.From, read.Cmd, read.Payload)
		}
	}
}

func main() {
	client, err := codechat.Connect("username")
	
	if err != nil {
		log.Println("could not connect")
		return
	}
	
	defer client.Conn.Close()
	
	go runReader(client)
	
	for {
		read := make([]byte, 4096)
		n, err := os.Stdin.Read(read)
		if err != nil || n == 0 {
			log.Println(err)
			return
		}
		readStr := strings.TrimSpace(string(read))
		client.Write("msg", readStr)
		//fmt.Println(b)
	}
	
}
