#!/usr/bin/env python2

import sys
import json
import pygtk
import gtksourceview2

import pango
import gobject
import gtk

def win_quit():
	print("peace out!")
	gtk.MainQuit()

class Layout:

	def __init__(self):
		win = gtk.Window()
		win.set_position(gtk.WIN_POS_CENTER)
		win.set_title("CodeChat")
		win.set_icon_name("gtk-dialog-info")
		win.connect("destroy", win_quit)

		leftEvent = gtk.EventBox()
		rightEvent = gtk.EventBox()
		#// main horizontal container, equal 1px spacing
		mainBox = gtk.HBox(True, 1)

		leftFrame = gtk.Frame("editor")
		rightFrame = gtk.Frame("chat")

		#// left and right containers, spacing not equal
		leftVBox = gtk.VBox(False, 1)
		rightVBox = gtk.VBox(False, 1)

		#// build up the left side
		editWindow = gtk.ScrolledWindow(None,None)
		editWindow.set_policy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		editWindow.set_shadow_type(gtk.SHADOW_IN)

		self.EditBuffer = gtksourceview2.Buffer()
		self.EditLangMgr = gtksourceview2.LanguageManager()
		self.EditBuffer.set_highlight_syntax(True)
		self.EditBuffer.set_highlight_matching_brackets(True)
		lang = self.EditLangMgr.get_language("python")
		self.EditBuffer.set_language(lang)

		self.EditView = gtksourceview2.View(self.EditBuffer)
		self.EditView.set_highlight_current_line(True)
		self.EditView.set_show_line_numbers(True)
		self.EditView.modify_font(pango.FontDescription("Monospace 8"))
		self.EditStatusBar = gtk.Statusbar()
		editWindow.add(self.EditView)
		
		leftVBox.pack_start(editWindow, True, True, 1)
		leftVBox.pack_end(self.EditStatusBar, False, False, 2)
		
		leftFrame.add(leftVBox)

		#// build up the right side
		chatWindow = gtk.ScrolledWindow(None, None)
		chatWindow.set_policy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		chatWindow.set_shadow_type(gtk.SHADOW_IN)
		
		chatTagTable = gtk.TextTagTable()
		self.ChatBuffer = gtk.TextBuffer(chatTagTable)
		self.ChatView = gtk.TextView(self.ChatBuffer)
		self.ChatView.set_editable(False)
		self.ChatView.set_cursor_visible(False)
		chatWindow.add(self.ChatView)
		chatHBox = gtk.HBox(False, 0)
		self.MessageBuffer = gtk.EntryBuffer("", -1)
		self.ChatEntry = gtk.Entry()
		self.ChatEntry.set_buffer(self.MessageBuffer)
		self.SendBtn = gtk.Button("send")

		chatHBox.pack_start(self.ChatEntry, True, True, 0)
		chatHBox.pack_end(self.SendBtn, False, False, 0)

		rightVBox.pack_start(chatWindow, True, True, 1)
		rightVBox.pack_start(chatHBox, False, False, 0)
		
		rightFrame.add(rightVBox)
		
		leftEvent.add(leftFrame)
		rightEvent.add(rightFrame)
		
		#// set up interface + layout
		mainBox.pack_start(leftEvent, True, True, 0)
		mainBox.pack_start(rightEvent, True, True, 0)
		
		#// Show the window
		win.add(mainBox)
		win.set_size_request(800, 400)
		win.show_all()

	def messageAction():
		print("message Action:")

def main():
	lyt = Layout()

	lyt.SendBtn.connect("clicked", lyt.messageAction)

	#lyt.ChatEntry.Connect("activate", func () {
		#messageAction(client, lyt)
	#})

	#// User has entered some text in the editor window
	#// send update-file 
	#lyt.EditBuffer.Connect("end-user-action", func() {
		#WriteLock.Lock()
		#var start, end gtk.TextIter
		#ctx := lyt.EditStatusBar.GetContextId("CodeChat")
		#lyt.EditStatusBar.Pop(ctx)
		#lyt.EditStatusBar.Push(ctx, "last edited by you")
		#lyt.EditBuffer.GetStartIter(&start)
		#lyt.EditBuffer.GetEndIter(&end)
		#file := lyt.EditBuffer.GetText(&start, &end, true)
		#client.Write("update-file",file)
		#WriteLock.Unlock()
	#})

	# need to:
	# [] connect to the server
	# [] connect event handlers
	# [] create startup dialog
	# [] quit gracefully
	gtk.main()

if __name__ == '__main__':
	main()
