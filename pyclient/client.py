'''
 client.py
 Authors: Graham Greving, David Taylor, Jake VanAdrighem
 CMPS 112: Final Project - CodeChat
 This program slowly evolved from a simple way to test the server
 to a fully fledged client program, and then into a library/package
 which is importable by any other python program. 

 It also exports the following functions:
 		ConnectionClient.Read()
                ConnectionClient.Write()
 		ConnectionClient.Connect()
                ConnectionClient.Close()
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
        try:
            ret = self.con.connect_ex((self.ip,self.port))
        except:
            #return 1 to indicate failure
            return 1
        return ret

    def Close(self):
        self.con.close()
        
    def Read(self):
        if not self.con:
            return()
        data = self.con.recv()
        if json.dumps(data):
            return data


    def Write(self, cmd, data):
        if not self.con or not data:
            return()
        sendMsg = json.dumps({'cmd':cmd, 'payload':data, 'from':self.username})
        self.con.sendall(sendMsg)
