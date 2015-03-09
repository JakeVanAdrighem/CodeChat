package main

import (
	codechat "CodeChat/client"
	"CodeChat/layout"
	"github.com/mattn/go-gtk/gtk"
	"log"
)

func doRead(client *codechat.Client, lyt *layout.Layout) {
	buffer := lyt.ChatBuffer
	var end gtk.TextIter
	for {
		read, err := client.Read()
		if err != nil {
			log.Println(err)
			if err.Error() == "EOF" {
				return
			}
		}
		switch read.Cmd {
		case "success":
			log.Println("success")
		case "message":
			buffer.GetEndIter(&end)
			buffer.Insert(&end, read.From + ": " + read.Payload +"\n")
		case "client-exit":
			buffer.GetEndIter(&end)
			buffer.Insert(&end, read.From + " has quit (" + read.Payload + ")\n")
		case "client-connect":
			buffer.GetEndIter(&end)
			buffer.Insert(&end, read.From + " has entered.\n")
		default:
			log.Println(read.From, read.Cmd, read.Payload)
		}
	}
}

func messageAction(client *codechat.Client, lyt *layout.Layout) {
	msg := lyt.MessageBuffer.GetText()
	println("button clicked: ", msg)
	var end gtk.TextIter
	buffer := lyt.ChatBuffer
	buffer.GetEndIter(&end)
	buffer.Insert(&end, "you: "+msg+"\n")
	if s:=len(msg); s != 0 {
		lyt.MessageBuffer.DeleteText(0, s)
	}
	// write msg here
	client.Write("msg", msg)
}

func main() {
	var name string
	var err error
	var client *codechat.Client

	lyt := new(layout.Layout)

	lyt.Init()

	// Send a message to the server when the input button is clicked
	lyt.SendBtn.Clicked(func () {
		messageAction(client, lyt)
	})
	// Send a message to the server when the user hits enter in the text
	// input
	lyt.ChatEntry.Connect("activate", func () {
		messageAction(client, lyt)
	})

	//// User has entered some text in the editor window
	//// send update-file 
	//layout.editorBuf.Connect("end-user-action", func() {
		//var start, end gtk.TextIter
		//layout.editorBuf.GetStartIter(&start)
		//layout.editorBuf.GetEndIter(&end)
		//log.Println("New file:")
		//log.Println(layout.editorBuf.GetText(&start, &end, true))
	//})


	// When focus enters the right side (editor):
	//		- set editor uneditable
	//		- set client.writeaccess false
	//		- send yield-write-access command
	// When focus enters the left side (chatmsgs or chat entry):
	//		- send request-write-access command
	//		- test client.writeaccess
	//		- if true -> set left editable and let the edits flow
	//		- if false -> wait (..?)



	// Connect the client
	name = layout.PromptUsername()
	client, err = codechat.Connect(name)
	if err != nil {
		log.Println("could not connect")
		return
	} else {
		var end gtk.TextIter
		lyt.ChatBuffer.GetEndIter(&end)
		lyt.ChatBuffer.Insert(&end, "Successfully connected as " + client.Username + "\n")
	}

	defer client.Close("F THIS")

	go doRead(client, lyt)
	gtk.Main()
}
