package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"unsafe"
)

const (
	HOST       = "localhost"
	PORT       = "8080"
	TYPE       = "tcp"
	BUFFERSIZE = 1024
	PATH       = "datafile.txt"
)

type IPApiData struct {
	CountryCode string
	Region      string
	Query       string
}

var (
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	lines  = []string{}
)

func GetExtData() IPApiData {
	request, err := http.Get("http://ip-api.com/json/")
	handleError(err)
	defer request.Body.Close()

	body, err := io.ReadAll(request.Body)
	handleError(err)

	var extData IPApiData
	json.Unmarshal(body, &extData)

	return extData
}

// Fill string with : to fit into buffer
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
	// Connect to server
	CONNECT := HOST + ":" + PORT
	conn, err := net.Dial(TYPE, CONNECT)
	handleError(err)
	defer conn.Close()

	sensitiveData := GetExtData()
	lines = append(lines,
		"IP: "+sensitiveData.Query,
		"Country: "+sensitiveData.CountryCode,
		"Region: "+sensitiveData.Region)

	// Transfer file with collected data
	SendData(conn)
}

func SendData(conn net.Conn) {
	arraySize := fillString(strconv.FormatInt(int64(unsafe.Sizeof(lines[0])*uintptr(len(lines))), 10), 10)
	conn.Write([]byte(arraySize))

	for i := 0; i < len(lines); i++ {
		sendBuffer := []byte(lines[i] + "\n")
		conn.Write(sendBuffer)
	}
}

func handleError(err error) {
	if err != nil {
		logger.Println("Error: ", err)
		return
	}
}
