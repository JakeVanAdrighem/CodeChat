package main

import (
    "log"
    "encoding/json"
    "os"
    "fmt"
    codechat "CodeChat/goclient/client"
)

var server string
var name string
var gochan chan string
var done int

func doReadTest(c *codechat.Client) {
    var v map[string]interface{}
    for {
        err := c.Jreader.Decode(&v)
        if err != nil {
            gochan <- c.Username + "read failed" 
        }
    }
}

func doWriteTest(c *codechat.Client) {
    var i int
    m := codechat.WriteMessage{"msg","stress-test-msg"}
    b, _ := json.Marshal(m)
    var v map[string]interface{}
    for i = 0; i < 10000; i++ {
        //client.Write("msg", "stress testing")
        c.Conn.Write(b)
        err := c.Jreader.Decode(&v)
        if err != nil {
            break
        }
    }
    gochan <- c.Username + "done" 
}

func main() {
    if len(os.Args) != 3 {
        log.Fatal("need username + server")
    } else {
        name = os.Args[1]
        server = os.Args[2]
    }
    gochan = make(chan string)
    for i := 0; i < 5; i++ {
        s := fmt.Sprintf("%s%d", name, i)
        c,err := codechat.Connect(s, server)
        defer c.Close("bye")
        if err != nil {
            log.Fatal("can't connect.")
        }
        go doWriteTest(c)
        //go doReadTest(c)
    }
    done := 0
    for x := range gochan {
        log.Println(x)
        done++
        if done == 5 {
            break
        }
    }
    log.Println("done")
}

