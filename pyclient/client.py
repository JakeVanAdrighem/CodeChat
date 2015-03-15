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
            sendMsg = json.dumps({'cmd':"connect", 'username':self.username})
            self.con.sendall(sendMsg)
        except:
            #return 1 to indicate failure
            return 1
        return ret

    def Close(self, reason):
        self.Write("exit", reason)
        self.con.close()
        
    def Read(self):
        try:
            data = self.con.recv(4096)
            x = json.loads(data.decode('ascii')) 
        except:
            print("no recv/loads")
            return
        if x:
            return x
        else:
            print("bad data")

    def Write(self, cmd, data):
        if not self.con or not data:
            return
        print("ConnectionClient.Write: " + self.username + "," + data + "," + cmd + ".")
        sendMsg = json.dumps({'cmd':cmd, 'msg':data})
        print(sendMsg)
        self.con.sendall(sendMsg)
