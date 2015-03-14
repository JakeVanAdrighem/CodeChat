package main

import (
    "log"
    "encoding/json"
    "os"
    codechat "CodeChat/goclient/client"
)

func main() {
    if len(os.Args) != 3 {
        log.Fatal("need username + server")
    }
    log.Println("start")
    client,err := codechat.Connect(os.Args[1],os.Args[2])
    log.Println("connected")
    defer client.Close("bye")

    if err != nil {
        log.Fatal("can't connect.")
    }
    var v map[string]interface{}
    var i int
    m := codechat.WriteMessage{"msg","stress-test-msg"}
    b, _ := json.Marshal(m)
    for i = 0; i < 10000; i++ {
        //client.Write("msg", "stress testing")
        client.Conn.Write(b)
        err := client.Jreader.Decode(&v)
        if err != nil {
            break
        }
        log.Println(i)
    }
}

