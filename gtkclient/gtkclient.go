package main

import (
	//	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	// "os"
	// "os/exec"
	// "path"
	// "regexp"
	// "sort"
	// "strings"
)

func main() {
	//var menuitem *gtk.MenuItem
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("CodeChat")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "CodeChat")

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

	// layout frames
	mainFrame := gtk.NewFrame("main")
	mainBox := gtk.NewHBox(true, 1)

	leftFrame := gtk.NewFrame("editor")
	//leftBox := gtk.NewVBox(true, 1)
	leftBox := gtk.NewScrolledWindow(nil, nil)
	leftBox.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	leftBox.SetShadowType(gtk.SHADOW_IN)
	editor := gtk.NewTextView()
	leftBox.Add(editor)

	leftFrame.Add(leftBox)

	rightFrame := gtk.NewFrame("chat")
	rightBox := gtk.NewVBox(false, 1)

	chatFrame := gtk.NewFrame("messages")
	//chatBox := gtk.NewVBox(true, 1)
	chatBox := gtk.NewScrolledWindow(nil, nil)
	chatBox.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	chatBox.SetShadowType(gtk.SHADOW_IN)
	chatMessages := gtk.NewTextView()
	chatMessages.SetEditable(false)
	chatMessages.SetCursorVisible(false)
	chatBox.Add(chatMessages)
	chatFrame.Add(chatBox)

	inputFrame := gtk.NewFrame("input")
	inputBox := gtk.NewHBox(false, 1)
	inputEntry := gtk.NewEntry()
	inputButton := gtk.NewButtonWithLabel("send")
	inputBox.Add(inputEntry)
	inputBox.Add(inputButton)
	inputFrame.Add(inputBox)

	rightBox.Add(chatFrame)
	rightBox.Add(inputFrame)

	rightFrame.Add(rightBox)

	mainBox.Add(leftFrame)
	mainBox.Add(rightFrame)
	mainFrame.Add(mainBox)

	window.Add(mainFrame)
	window.SetSizeRequest(1000, 600)
	window.ShowAll()
	gtk.Main()
}
