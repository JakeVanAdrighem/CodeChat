package main

import (
	codechat "CodeChat/client"
	"CodeChat/layout"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"log"
	"sync"
)

var WriteLock sync.Mutex

func doRead(client *codechat.Client, lyt *layout.Layout) {
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
			buffer := lyt.ChatBuffer
			var end gtk.TextIter
			buffer.GetEndIter(&end)
			buffer.Insert(&end, read.From + ": " + read.Payload +"\n")
		case "client-exit":
			buffer := lyt.ChatBuffer
			var end gtk.TextIter
			buffer.GetEndIter(&end)
			buffer.Insert(&end, read.From + " has quit (" + read.Payload + ")\n")
		case "client-connect":
			buffer := lyt.ChatBuffer
			var end gtk.TextIter
			buffer.GetEndIter(&end)
			buffer.Insert(&end, read.From + " has entered.\n")
		case "update-file":
			gdk.ThreadsEnter()
			//lyt.EditBuffer.SetEditable(false)
			//lyt.EditBuffer.SetCursorVisible(false)
			//lang := lyt.EditLangMgr.GuessLanguage("",read.Payload)
			//if *lang != nil {
				////log.Println(lang.GetStyleIds())
				//log.Println(lang)
				//lyt.EditBuffer.SetLanguage(lang)
			//}
			ctx := lyt.EditStatusBar.GetContextId("CodeChat")
			lyt.EditStatusBar.Pop(ctx)
			lyt.EditStatusBar.Push(ctx, "last edited by " + read.From)
			lyt.EditBuffer.SetText(read.Payload)
			gdk.ThreadsLeave()
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
	var name,ipport string
	var err error
	var client *codechat.Client

	lyt := new(layout.Layout)
	lyt.Init()

	// Send a message to the server when the input button is clicked
	lyt.SendBtn.Clicked(func () {
		messageAction(client, lyt)
	})
	// Send a message to the server when the user hits enter
	lyt.ChatEntry.Connect("activate", func () {
		messageAction(client, lyt)
	})
	// User has entered some text in the editor window
	// send update-file 
	lyt.EditBuffer.Connect("end-user-action", func() {
		WriteLock.Lock()
		var start, end gtk.TextIter
		ctx := lyt.EditStatusBar.GetContextId("CodeChat")
		lyt.EditStatusBar.Pop(ctx)
		lyt.EditStatusBar.Push(ctx, "last edited by you")
		lyt.EditBuffer.GetStartIter(&start)
		lyt.EditBuffer.GetEndIter(&end)
		file := lyt.EditBuffer.GetText(&start, &end, true)
		client.Write("update-file",file)
		WriteLock.Unlock()
	})

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
	name, ipport = layout.PromptUsername()
	if name == "" {
		name = "dumbass"
	}
	// connect locally if there is no 
	if ipport == "" {
		ipport = "127.0.0.1:8080"
	}
	client, err = codechat.Connect(name, ipport)
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
	gdk.ThreadsInit()
	gtk.Main()
}
