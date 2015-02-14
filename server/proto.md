CodeChat Client-Server Protocol
===============================
This document outlines the general protocol for communication between
CodeChat clients and servers. Communication and actions are facilitated
by a set of commands and arguments formated as JSON with the form:

	{
		"cmd":"<command>",
		[args ..]
	}

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
		"old-username":"<old-username>",
		"new-username":"<new-username>"
	}

### 3. exit()
Disconnects the client from the server:

	{
		"cmd":"exit"
	}

##Return Messages
Status of commands are always returned back to the client from the
server as JSON:

	{
		"success":<bool>,
		"error-msg": "<error>"
	}

Success will be true or false depending on success or failure of the previous command,
in the case of a failure, an error message is provided.
