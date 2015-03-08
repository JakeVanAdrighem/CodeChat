package main

import (
	codechat "CodeChat/client"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gtksourceview"
	"log"
)

type Layout struct {
	// frames + panes
	mainFrame  gtk.IWidget
	leftFrame  gtk.IWidget
	rightFrame gtk.IWidget

	// objects
	editor       *gtksourceview.SourceView
	inputEntry   *gtk.Entry
	inputButton  *gtk.Button
	chatMessages *gtk.TextView
}

/*
							main
	+----------------------------------------+
	|		left		|		right		|
	|					|					|
	|					|					|
	|					|					|
	|					|					|
	|					+-------------------+
	|					|				|btn|
	+--------------------+-------------------+

*/
func layoutInit() Layout {
	// layout frames
	mainFrame := gtk.NewFrame("")
	mainBox := gtk.NewHBox(true, 1)
	leftFrame := gtk.NewFrame("editor")
	leftBox := gtk.NewScrolledWindow(nil, nil)
	leftBox.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	leftBox.SetShadowType(gtk.SHADOW_IN)
	// need to figure out this sourceview thing
	editorbuf := gtksourceview.NewSourceBuffer()
	editor := gtksourceview.NewSourceViewWithBuffer(editorbuf)
	leftBox.Add(editor)
	leftFrame.Add(leftBox)
	rightPane := gtk.NewVPaned()
	chatFrame := gtk.NewFrame("chat")
	chatBox := gtk.NewScrolledWindow(nil, nil)
	chatBox.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	chatBox.SetShadowType(gtk.SHADOW_IN)
	chatMessages := gtk.NewTextView()
	chatMessages.SetEditable(false)
	chatMessages.SetCursorVisible(false)
	chatBox.Add(chatMessages)
	chatFrame.Add(chatBox)
	inputFrame := gtk.NewFrame("")
	inputBox := gtk.NewHBox(false, 1)
	inputEntry := gtk.NewEntry()
	inputButton := gtk.NewButtonWithLabel("send")
	inputBox.Add(inputEntry)
	inputBox.Add(inputButton)
	inputFrame.Add(inputBox)
	rightPane.Pack1(chatFrame, false, false)
	rightPane.Pack2(inputFrame, false, false)
	mainBox.Add(leftFrame)
	mainBox.Add(rightPane)
	mainFrame.Add(mainBox)
	inputEntry.SetSizeRequest(450, -1)
	chatFrame.SetSizeRequest(500, 550)
	inputFrame.SetSizeRequest(500, 50)

	return Layout{mainFrame, leftFrame, rightPane, editor, inputEntry, inputButton, chatMessages}
}

type Message struct {
	Cmd string `json:"cmd"`
	Msg string `json:"msg"`
}

type Connect struct {
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
}

func doRead(client *codechat.Client, lyt *Layout) {
	buffer := lyt.chatMessages.GetBuffer()
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

func main() {
	//var menuitem *gtk.MenuItem

	var name string
	var err error
	var client *codechat.Client

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("CodeChat")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "CodeChat")

	layout := layoutInit()

	// Send a message to the server when the input button is clicked
	layout.inputButton.Clicked(func() {
		// send mesage here
		msg := layout.inputEntry.GetText()
		println("button clicked: ", msg)
		var end gtk.TextIter
		buffer := layout.chatMessages.GetBuffer()
		buffer.GetEndIter(&end)
		buffer.Insert(&end, "you: "+msg+"\n")
		layout.inputEntry.SetText("")
		// write msg here
		client.Write("msg", msg)
	})
	// Send a message to the server when the user hits enter in the text
	// input
	layout.inputEntry.Connect("activate", func() {
		msg := layout.inputEntry.GetText()
		println("enter pressed: ", msg)
		var end gtk.TextIter
		buffer := layout.chatMessages.GetBuffer()
		buffer.GetEndIter(&end)
		buffer.Insert(&end, "you: "+msg+"\n")
		layout.inputEntry.SetText("")
		// write msg here
		client.Write("msg", msg)
	})

	// Show a dialog to get the username on startup
	messagedialog := gtk.NewDialog()
	connectBox := messagedialog.GetVBox()
	label := gtk.NewLabel("username")
	username := gtk.NewEntry()
	username.Connect("activate", func() {
		name = username.GetText()
		messagedialog.Destroy()
	})
	connectBox.Add(label)
	connectBox.Add(username)
	messagedialog.AddButton("connect", gtk.RESPONSE_OK)
	messagedialog.Response(func() {
		name = username.GetText()
	})
	label.Show()
	username.Show()
	messagedialog.Run()
	messagedialog.Destroy()
	// End Dialog
	
	// Setup the window
	window.Add(layout.mainFrame)
	window.SetSizeRequest(1000, 600)
	window.ShowAll()

	// Connect the client
	client, err = codechat.Connect(name)
	
	if err != nil {
		log.Println("could not connect")
		return
	} else {
		var end gtk.TextIter
		buffer := layout.chatMessages.GetBuffer()
		buffer.GetEndIter(&end)
		buffer.Insert(&end, "Successfully connected as " + client.Username + "\n")
	}

	defer client.Close("F THIS")

	go doRead(client, &layout)

	gtk.Main()
}
