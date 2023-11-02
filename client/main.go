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
	PATH       = "dummyfile.txt"
	BUFFERSIZE = 1024
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

// Fill string with bytes
func fillString(retunString string, toLength int) string {
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
	CONNECT := HOST + ":" + PORT
	conn, err := net.Dial(TYPE, CONNECT)
	handleError(err)
	defer conn.Close()

	SendFile(conn, PATH)
}

func SendFile(conn net.Conn, path string) {
	payload, err := os.Open(path)
	handleError(err)
	fileInfo, err := payload.Stat()
	handleError(err)
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	sendBuffer := make([]byte, BUFFERSIZE)
	conn.Write([]byte(fileSize))
	for {
		_, err = payload.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
}

func handleError(err error) {
	if err != nil {
		logger.Println("Error: ", err)
		return
	}
}
