# FileDriver
An easy way to save files in your server.
**FileDriver** is a CLI tool written in Go to manage files from a server.

** IMPORTANT: This project still in development but is working with the most important functionalities. **

## Installation
For now you only can build the project by your own and execute it.

### Requirements
- **Go**: Is necesary to have installed Go **v1.21 or superior**. You can download it from [golang.org](https://golang.org/dl/).

### Client installation
For the client you can run the command to install the executable.
    ```
	$ go install github.com/alexismrosales/FileDriver/cmd/filedriver@latest
    ```

### Server installation
To install the server run:
    ```
	$ go install github.com/alexismrosales/FileDriver/cmd/filedriverserver@latest
    ```

## Quickstart

### Client user
Running the command:
    ```
    $ filedriver --help
    ```
You can read all avalaible commands.

To create a connection with the server you can run once:
    ```
    $ filedriver setaddr ip port
    ```
the address will be save, so do it once is enough.



    ```
    $ filedriver setaddr ip port
    ```
the address will be save, so do it once is enough.
### Server
Running
You can check this page here: https://alexismrosales.github.io/
