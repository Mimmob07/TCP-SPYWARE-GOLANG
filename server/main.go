package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	HOST       = "localhost"
	PORT       = "8080"
	TYPE       = "tcp"
	BUFFERSIZE = 1024
	OUTPUTPATH = "RecievedData/"
)

var count = 0
var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

func handleError(err error) {
	if err != nil {
		logger.Println("Panic: ", err)
		return
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	fileSizeBuffer := make([]byte, 10)
	conn.Read(fileSizeBuffer)
	fileSize, err := strconv.ParseInt(strings.ReplaceAll(string(fileSizeBuffer), ":", ""), 10, 64)
	handleError(err)
	logger.Printf("File being recived of size %v bytes\n", fileSize)

	var fileName string = OUTPUTPATH + "data" + fmt.Sprint(count) + ".txt"

	createPayload, err := os.Create(fileName)
	handleError(err)
	defer createPayload.Close()

	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(createPayload, conn, (fileSize - receivedBytes))
			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(createPayload, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	logger.Printf("%v has been created\n", fileName)
}

func main() {
	server, err := net.Listen(TYPE, HOST+":"+PORT)
	handleError(err)
	defer server.Close()
	logger.Println("Server is now listening")

	for {
		conn, err := server.Accept()
		logger.Println("New connection accepted")
		handleError(err)
		go HandleConnection(conn)
		count++
	}
}
