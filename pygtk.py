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

def layout():
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

	EditBuffer = gtksourceview2.Buffer()
	EditLangMgr = gtksourceview2.LanguageManager()
	EditBuffer.set_highlight_syntax(True)
	EditBuffer.set_highlight_matching_brackets(True)
	lang = EditLangMgr.get_language("python")
	EditBuffer.set_language(lang)

	EditView = gtksourceview2.View(EditBuffer)
	EditView.set_highlight_current_line(True)
	EditView.set_show_line_numbers(True)
	EditView.modify_font(pango.FontDescription("Monospace 8"))
	EditStatusBar = gtk.Statusbar()
	editWindow.add(EditView)
	
	leftVBox.pack_start(editWindow, True, True, 1)
	leftVBox.pack_end(EditStatusBar, False, False, 2)
	
	leftFrame.add(leftVBox)

	#// build up the right side
	chatWindow = gtk.ScrolledWindow(None, None)
	chatWindow.set_policy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	chatWindow.set_shadow_type(gtk.SHADOW_IN)
	
	chatTagTable = gtk.TextTagTable()
	ChatBuffer = gtk.TextBuffer(chatTagTable)
	ChatView = gtk.TextView(ChatBuffer)
	ChatView.set_editable(False)
	ChatView.set_cursor_visible(False)
	chatWindow.add(ChatView)
	chatHBox = gtk.HBox(False, 0)
	MessageBuffer = gtk.EntryBuffer("", -1)
	ChatEntry = gtk.Entry()
	ChatEntry.set_buffer(MessageBuffer)
	SendBtn = gtk.Button("send")

	chatHBox.pack_start(ChatEntry, True, True, 0)
	chatHBox.pack_end(SendBtn, False, False, 0)

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

def main():
	layout()
	# need to:
	# [] connect to the server
	# [] connect event handlers
	# [] create startup dialog
	# [] quit gracefully
	gtk.main()

if __name__ == '__main__':
	main()
