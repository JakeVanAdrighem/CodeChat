#!/usr/bin/env python2

import sys
import json
import pygtk
import gtksourceview2

import pango
import gobject
import gtk

import threading
import time

import client

def win_quit(whatisthis):
    print("peace out!")
    gtk.main_quit()

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

        self.conclient = client.ConnectionClient()

        #we need to query the user for username and connection info
        self.dialog = gtk.Dialog("Connection Dialog",win,gtk.DIALOG_MODAL)
        ulabel = gtk.Label("username")
        ilabel = gtk.Label("IP Address:Port")
        self.uentry = gtk.Entry()
        self.ientry = gtk.Entry()
        self.dialog.vbox.add(ulabel)
        self.dialog.vbox.add(self.uentry)
        self.dialog.vbox.add(ilabel)
        self.dialog.vbox.add(self.ientry)
        self.button = self.dialog.add_button("connect",gtk.RESPONSE_OK)
        self.button.connect("clicked",self.connect)
        self.ientry.connect("activate",self.connect)
        #self.dialog.response(self.connect)
        self.dialog.show_all()
        response = self.dialog.run()
        print("here- after dialog")
        #instantiate read thread
        self.read_thread = threading.Thread(None, self.doRead, 'read_thread')
        self.read_thread.start()


    def doRead(self):
        #eventually implement it so that we
        #update less frequently when we're editing
        #if not self.editing:
        print( "in doRead thread")
        while True:
            print("here")
            res = self.conclient.Read()
            print ("got message" + json.dumps(res))
            if res:
                #gtk.gdk.threads_enter()
                try:
                    cmd = res["cmd"]
                except:
                    cmd = "none"
                if cmd == "return-status":
                    print("return status")
                if cmd == "update-file":
                    ctx = self.EditStatusBar.get_context_id("CodeChat")
                    self.EditStatusBar.pop(ctx)
                    self.EditStatusBar.push(ctx, "last edited by " + res["from"])
                    self.EditBuffer.SetText(res["payload"])
                elif cmd   == "client-connect":
                    endIter = self.ChatBuffer.get_end_iter()
                    self.ChatBuffer.insert(endIter,res["from"] + " has entered.\n")
                #elif cmd == "success":
                    #successful connection
                elif cmd == "client-exit":
                    endIter = self.ChatBuffer.get_end_iter()
                    self.ChatBuffer.insert(endIter,res["from"] + " has quit (" + res["payload"] + ")\n")
                elif cmd == "none":
                    print("weird shit happened")
                #pause thread for a quarter second
                #time.sleep(0.25)
                #gtk.gdk.threads_exit()

    def connect(self, whatisthis):
        #we have to include the username as a member
        #only because we post messages locally and
        #need our own name
        self.username = self.uentry.get_text()
        ip       = self.ientry.get_text()
        #don't destroy the dialog if they don't provide proper connection info
        try:
                ip,port = ip.split(':')
        except:
                print("IP and Port bad formatting or not provided")
                return
        res = self.conclient.Connect(self.username,ip,port)
        self.conclient.Write("connect","username")
        #res will be 0 if it's a successful connection
        if res:
                return
        self.dialog.destroy()

    def messageAction(self, whatisthis):
            input = self.ChatEntry.get_text()
            if input in ('','\n'):
                return
            self.conclient.Write("msg",input)
            print("message Action:" + input)
            self.ChatEntry.set_text('')
            self.ChatBuffer.insert(self.ChatBuffer.get_end_iter(),self.username + ": " + input + '\n')

    def editorAction(self, whatisthis):
        ctx = self.EditStatusBar.get_context_id("CodeChat")
        self.EditStatusBar.pop(ctx)
        self.EditStatusBar.push(ctx, "last edited by you")
        start = self.EditStatusBar.get_start_iter()
        end   = self.EditStatusBar.get_end_iter()
        data  = self.EditBuffer.get_text(start, end, True)
        self.conclient.Write("update-file", data)

def main():
    lyt = Layout()
    lyt.SendBtn.connect("clicked", lyt.messageAction)
    lyt.ChatEntry.connect("activate", lyt.messageAction)
    gtk.main()

if __name__ == '__main__':
    main()
