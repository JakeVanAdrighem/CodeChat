# CodeChat #
Collaborative Development Environment

CMPS 112

2015

By: Graham Greving, David Taylor, and Jake VanAdrighem

## Overview ##
CodeChat is a collaborative code development environment which
combines the features of a Google Docs style editior with features
from a source-code editor, and an IM service. It serves to facilitate
remote collaborative development.

The project was implemented entirely in Go, and includes
modules for both the CodeChat server and client. An additional
Python client is provided to demonstrate the server's ability
to handle clients writen in multiple languages.

## Build Instructions ##
As of this edit, no binaries are pre-bundled, and all
builds are from source.

Most of the tools used are in the Go standard library,
and builds are fairly easy.

Be sure to have all dependencies installed prior to building.

### Dependencies ###
* go
* go-gtk: <github.com/gramasaurous/go-gtk>
* gtk2
* gtk-sourceview

Additionally, for the python client:
* python2
* pygtk
* pygtksourceview
* gobject

To build the server:

    $ cd server
    $ go build

To build the client:

    $ cd goclient
    $ go build

## Running ##

In order to use the client, you must first set up the server:

    $ cd server
    $ ./server <port>
    # where <port> is the port you want to run on

Next you can fire up the go client:
    
    $ cd goclient
    $ ./goclient

This will prompt you with a dialog box asking for a username and
an address. The address should be of the form "address:port".
If blank, address will default to "localhost:8080".

Similarly you can run the python client:

    $ cd pyclient
    $ python2 pyclient.py

The interface works the same across the two clients.

