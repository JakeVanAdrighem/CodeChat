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
		"username":"<username"
	}
### 2. rename ("old-username", "new-username")
Changes the current client's username:

	{
		"cmd":"rename",
		"oldname":"<old-username>",
		"newname":"<new-username>"
	}

### 3. get-clients()

	{
		"cmd":"get-clients"
	}

### 4. exit()
Disconnects the client from the server:

	{
		"cmd":"exit"
	}

##Status Messages
Status of commands are always returned back to the client from the
server as JSON:

	{
		"success":<bool>,
		"error-msg": "<error>",
		"success-msg":"<success-message>"
	}

Success will be true or false depending on success or failure of the previous command, in the case of a failure, an error message is provided
in the "error-msg" field. In commands where information is requested from
the server, this info is passed back to the server in the "success-msg"
field.

##General Messages
Messages indicating server state are also passed back to clients for them
to be displayed. The following messages are sent to each client depending
on various server events:

When a new message has been sent to the chat:

	{
		"cmd":"message",
		"user":"username"
	}

When a new client enters the chat:

	{
		"cmd":"client-enter",
		"user":"<username>"
	}

When an existing client exits the chat:

	{
		"cmd":"client-exit",
		"user":"<username>"
	}

Note that clients are responsible for keeping a list of current clients
in the chat. So appropriate actions should be taken to notify the client
of entering and exiting users, as well as update a list of current users.
