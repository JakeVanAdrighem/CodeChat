'''
// client.py
// Authors: Graham Greving, David Taylor, Jake VanAdrighem
// CMPS 112: Final Project - CodeChat

// This program slowly evolved from a simple way to test the server
// to a fully fledged client program, and then into a library/package
// which is exportable to any other go program. It exports the structs:
// 		ReadMessage
//		WriteMessage
//		Client
//  (it technically exports the other ones, but only so that they can
//   be marshalled into json)
// It also exports the following functions:
// 		Client.Read()
// 		Client.Write()
// 		Connect
'''
import socket
import json



class ConnectionClient():
    def _init_(self):
        #build self
        self.username = 'NoName'
        self.ip = '127.0.0.1'
        self.port = '8080'

    def Connect(self, username, ip, port):
        #create connection
        self.username = username
        self.ip = ip
        #port must be int but is not assumed to be
        self.port = int(port)
        self.con = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.con.connect((self.ip,self.port))

    def Close(self):
        self.con.close()
        
    def Read(self):
        if not self.con:
            return()
        data = self.con.recv()
        if json.dumps(data):
            return data
        '''
        #check msg
        data = json.load(msg)
        command = data["cmd"]
        if command == "message":
            msg = data["payload"]
            msgSender = data["from"]
            #push message to chat
        elif command == "client-connect":
            #person has connected.
            user = data["payload"]
            #push connect message to chat
        elif command == "client-exit":
            #person has disconnected
            user = data["payload"]
            #push disconnect message to chat
        elif command == "update-file":
            fileData = data["payload"]
            #push file data to SourceView
        '''


    def Write(self, cmd, data):
        if not self.con or not data:
            return()
        sendMsg = json.dumps({'cmd':cmd, 'payload':data, 'from':self.username})
        self.con.sendall(sendMsg)
