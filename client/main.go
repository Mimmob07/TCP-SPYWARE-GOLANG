package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
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

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

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

// Bad way of doing this shit
func FormatAndWriteData(data IPApiData) {
	var (
		lines     = []string{"Country: " + data.CountryCode, "Region: " + data.Region, "IP: " + data.Query}
		file, err = os.Create(PATH)
	)
	handleError(err)
	defer file.Close()

	// dont need to do this, could just send the array of lines as bytes
	for i := 0; i < len(lines); i++ {
		_, err := file.WriteString(lines[i] + "\n")
		handleError(err)
		file.Sync()
	}
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
	FormatAndWriteData(sensitiveData)

	// Transfer file with collected data
	SendFile(conn, PATH)
}

func SendFile(conn net.Conn, path string) {
	// Load file and resources
	dataFile, err := os.Open(path)
	handleError(err)
	fileInfo, err := dataFile.Stat()
	handleError(err)
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)

	// Send file size
	sendBuffer := make([]byte, BUFFERSIZE)
	conn.Write([]byte(fileSize))

	// Send actual file
	for {
		_, err = dataFile.Read(sendBuffer)
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
