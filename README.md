# CodeChat #
Collaborative Development Environment
CMPS 112
2015

By: Graham Greving, David Taylor, and Jake VanAdrighem

## Overview ##
CodeChat is a collaborative code development environment which
combines the features of a Google Docs style editor with an
IM service, to facilitate remote collaborative development.

The project was implemented entirely in Go, and includes
modules for both the CodeChat server and client.

## Build Instructions ##
As of this edit, no binaries are pre-bundled, and all
builds are from source.

Most of the tools used are in the Go standard library,
and builds are fairly easy.

To build the server:

``
cd server
go build
``

To build the client:

``
cd gtkclient
go build
``

Be sure to have all dependencies installed prior to building.

### Dependencies ###
* go-gtk: <github.com/gramasaurous/go-gtk>
* gtk2
* gtk-sourceview
