package layout

import (
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/glib"
	gtksource "github.com/mattn/go-gtk/gtksourceview"
	//"log"
)
/*
 * 						mainHBox
 * 		leftFrame					rightFrame
 * 			|							|
 * 		leftEvent					rightEvent
 * 			|							|
 * 		leftVBox					leftVBox
 * 	|				|			  |				|
 * scrolledwindow   |		scrolledwindow  chatHBox
 *	|				|			 |			|		|
 * editorEntry  editorStatus 	chat	msgentry	sendbtn
 * 	
 */

type Layout struct {
	// window + main interfaces
	win *gtk.Window
	// main horizontal box
	mainBox *gtk.HBox
	// dividing frames
	leftFrame *gtk.Frame
	rightFrame *gtk.Frame
	// event boxes
	leftEvent *gtk.EventBox
	rightEvent *gtk.EventBox
	// left and right vertical boxes
	leftVBox *gtk.VBox
	rightVBox *gtk.VBox
	// left widgets
	editWindow	*gtk.ScrolledWindow
	EditBuffer  *gtksource.SourceBuffer
	editView	   *gtksource.SourceView
	editStatusBar  *gtk.Statusbar
	// right widgets
	chatWindow *gtk.ScrolledWindow
	chatTagTable *gtk.TextTagTable
	ChatBuffer *gtk.TextBuffer
	chatView   *gtk.TextView
	chatHBox   *gtk.HBox
	MessageBuffer *gtk.EntryBuffer
	ChatEntry  *gtk.Entry
	SendBtn    *gtk.Button
}

func (lyt *Layout) Init() {
	// set up window
	gtk.Init(nil)
	lyt.win = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	lyt.win.SetPosition(gtk.WIN_POS_CENTER)
	lyt.win.SetTitle("CodeChat")
	lyt.win.SetIconName("gtk-dialog-info")
	lyt.win.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("peace out!", ctx.Data().(string))
		gtk.MainQuit()
	}, "CodeChat")
	
	lyt.leftEvent = gtk.NewEventBox()
	lyt.rightEvent = gtk.NewEventBox()
	// main horizontal container, equal 1px spacing
	lyt.mainBox = gtk.NewHBox(true, 1)
	
	lyt.leftFrame = gtk.NewFrame("editor")
	lyt.rightFrame = gtk.NewFrame("chat")
	
	// left and right containers, spacing not equal
	lyt.leftVBox = gtk.NewVBox(false, 1)
	lyt.rightVBox = gtk.NewVBox(false, 1)
	
	// build up the left side
	lyt.editWindow = gtk.NewScrolledWindow(nil,nil)
	lyt.editWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	lyt.editWindow.SetShadowType(gtk.SHADOW_IN)
	
	lyt.EditBuffer = gtksource.NewSourceBuffer()
	lyt.editView = gtksource.NewSourceViewWithBuffer(lyt.EditBuffer)
	lyt.editStatusBar = gtk.NewStatusbar()
	lyt.editWindow.Add(lyt.editView)
	
	lyt.leftVBox.PackStart(lyt.editWindow, true, true, 1)
	lyt.leftVBox.PackEnd(lyt.editStatusBar, false, false, 2)
	
	lyt.leftFrame.Add(lyt.leftVBox)

	// build up the right side
	lyt.chatWindow = gtk.NewScrolledWindow(nil,nil)
	lyt.chatWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	lyt.chatWindow.SetShadowType(gtk.SHADOW_IN)
	
	lyt.chatTagTable = gtk.NewTextTagTable()
	lyt.ChatBuffer = gtk.NewTextBuffer(lyt.chatTagTable)
	lyt.chatView = gtk.NewTextViewWithBuffer(*lyt.ChatBuffer)
	lyt.chatView.SetEditable(false)
	lyt.chatView.SetCursorVisible(false)
	lyt.chatWindow.Add(lyt.chatView)
	lyt.chatHBox = gtk.NewHBox(false, 0)
	lyt.MessageBuffer = gtk.NewEntryBuffer("")
	lyt.ChatEntry = gtk.NewEntryWithBuffer(lyt.MessageBuffer)
	lyt.SendBtn = gtk.NewButtonWithLabel("send")
	
	lyt.chatHBox.PackStart(lyt.ChatEntry, true, true, 0)
	lyt.chatHBox.PackEnd(lyt.SendBtn, false, false, 0)
	
	lyt.rightVBox.PackStart(lyt.chatWindow, true, true, 1)
	lyt.rightVBox.PackStart(lyt.chatHBox, false, false, 0)
	
	lyt.rightFrame.Add(lyt.rightVBox)
	
	lyt.leftEvent.Add(lyt.leftFrame)
	lyt.rightEvent.Add(lyt.rightFrame)
	
	// set up interface + layout
	lyt.mainBox.PackStart(lyt.leftEvent, true, true, 0)	

	lyt.mainBox.PackStart(lyt.rightEvent, true, true, 0)
	
	// Show the window
	lyt.win.Add(lyt.mainBox)
	lyt.win.SetSizeRequest(800, 400)
	lyt.win.ShowAll()
}

func PromptUsername() (name string){
	// Show a dialog to get the username on startup
	dialog := gtk.NewDialog()
	connectBox := dialog.GetVBox()
	label := gtk.NewLabel("username")
	username := gtk.NewEntry()
	username.Connect("activate", func() {
		name = username.GetText()
		dialog.Destroy()
	})
	connectBox.Add(label)
	connectBox.Add(username)
	dialog.AddButton("connect", gtk.RESPONSE_OK)
	dialog.Response(func() {
		name = username.GetText()
		dialog.Destroy()
	})
	label.Show()
	username.Show()
	dialog.Run()
	//messagedialog.Destroy()
	// End Dialog
	return
}
