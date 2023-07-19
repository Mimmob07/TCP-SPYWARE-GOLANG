package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	HOST       = "localhost"
	PORT       = "8080"
	TYPE       = "tcp"
	BUFFERSIZE = 1024
)

var count = 0
var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

func handleError(err error) {
	if err != nil {
		logger.Println("Panic: ", err)
		return
	}
}

func handleConnection(conn net.Conn) {
	for {
		payload, err := os.Open("payload.txt")
		handleError(err)
		fileInfo, err := payload.Stat()
		handleError(err)
		defer conn.Close()
		fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
		fileName := fillString(fileInfo.Name(), 64)
		conn.Write([]byte(fileSize))
		conn.Write([]byte(fileName))
		sendBuffer := make([]byte, BUFFERSIZE)
		for {
			_, err = payload.Read(sendBuffer)
			if err == io.EOF {
				break
			}
			conn.Write(sendBuffer)
		}
		logger.Println("File sent in theory")
		conn.Close()
		break
	}
}

func fillString(retunString string, toLength int) string { // Fill string with bytes
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func main() {

	server, err := net.Listen(TYPE, HOST+":"+PORT)
	handleError(err)
	defer server.Close()
	logger.Println("Server is now listening")

	for {
		conn, err := server.Accept()
		handleError(err)
		go handleConnection(conn)
		count++
	}
}
