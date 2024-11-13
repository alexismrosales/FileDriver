# FileDriver
An easy way to save files in your server.
**FileDriver** is a CLI tool written in Go to manage files from a server.

**❗IMPORTANT: This project still in development but is working with the most important functionalities.**

## Installation
For now you only can build the project by your own and execute it.

### Requirements
- **Go**: Is necesary to have installed Go **v1.21 or superior**. You can download it from [golang.org](https://golang.org/dl/).

### Client installation
For the client you can run the command to install the executable.

	$ go install github.com/alexismrosales/FileDriver/cmd/filedriver@latest    
### Server installation
Exec this command to install the server.

	$ go install github.com/alexismrosales/FileDriver/cmd/fdserver@latest

## Quickstart

### Client user
Running the command:

    $ filedriver --help
    
You can read all avalaible commands.

To create a connection with the server you can run once:
    
    $ filedriver setaddr ip port
    
the address will be save, so do it once is enough.

You can download files using:

    $ filedriver download serverPaths
    
The paths to use are from the server directory, to see which files are availble use:

    $ filedriver ls

Also you can use a upload argument.

    $ filedriver upload paths

where you can use as any paths you want from your computer.    
### Server
To start listening petitions run:

    $ fdserver ip port
It is important to keep the application running to receive petitions from server.

### Example
As the client, first time connecting with server:

    $ filedriver setaddr 127.0.0.1 8080
will use the address *127.0.0.1:8080*.

To upload a file you run:

    $ filedriver upload file1 file2 dir1/ dir2/

you can also use shorted paths like this:

    $ filedriver upload ~/dir1/file1

**NOTE: There will be probably some problems, so use it as your own risk.**
