CodeChat Client-Server Protocol
===============================
This document outlines the general protocol for communication between
CodeChat clients and servers. Communication and actions are facilitated
by a set of commands and arguments formated as JSON with the form:

	{
		"cmd":"<command>",
		args ...
	}

At the moment this is a rudimentary protocol specification required for
minimal a chat-service. Commands will be expanded to support various
other operations, including, but not limited to: text-editing, chatrooms,
and authentication.

##Available commands are:
### 1. connect("username")
Adds the client to the chat server with the provided username:

	{
		"cmd":"connect",
		"username":"<username>"
	}

### 1. message ("message")

	{
		"cmd":"msg",
		"msg":<message>
	}

### 2. rename ("old-username", "new-username") ** not implemented **
Changes the current client's username:

	{
		"cmd":"rename",
//		"oldname":"<old-username>",
		"newname":"<new-username>"
	}

### 3. get-clients() ** not implemented **

	{
		"cmd":"get-clients"
	}

### 4. exit()
Disconnects the client from the server:

	{
		"cmd":"exit"
	}

## Text Editing Commands

### 1. requestWriteAccess() ** not implemented **

    {
        "cmd": "request-write-ccess"
    }

### 2. yieldWriteAccess() ** not implemented **

    {
        "cmd":"yield-write-access"
    }

### 3. updateFile(file)

    {
        "cmd":"update-file",
        "file":<file>
    }

##Status Messages
Status of commands are always returned back to the client from the
server as JSON:

	{
		"cmd":"return-status",
		"return-command":<command>,
		"status":<true,false>,
		"payload":<payload from the server, if any>
	}

Since all messages sent back from the server are indistinguishable,
return-status are formatted as a "return-status" command. The string
"return-command" contains the command that was recently attempted by the
server. "status" is a boolean on the success, and payload contains any
data sent back from the server (this is not always filled). An example
where "payload" is filled is after the completion of a "connect." In this
case, the payload contains the most recent version of the file being
edited in the chatroom. This ensures that newly connected clients will
have a copy of the file when they connect.

##General Messages
Messages indicating server state are also passed back to clients for them
to be displayed. The following messages are sent to each client depending
on various server events:

When a new message has been sent to the chat:

	{
		"cmd":"message",
		"payload":"message-from-client",
		"from":<from-client-username>
	}

When a new client enters the chat:

	{
		"cmd":"client-connect",
		"from":<from-client-username>,
		"payload":<username>
	}

When an existing client exits the chat:

	{
		"cmd":"client-exit",
		"from":<from-client-username>,
		"payload":"<username>"
	}

When the text file is updated:

    {
        "cmd":"update-file",
     	"from":<from-client-username>,
        "payload":<new-file>
    }

Note that clients are responsible for keeping a list of current clients
in the chat. So appropriate actions should be taken to notify the client
of entering and exiting users, as well as update a list of current users.
