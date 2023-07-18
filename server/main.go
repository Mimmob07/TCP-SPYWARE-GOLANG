package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
        listen, err := net.Listen(TYPE, HOST+":"+PORT)
        handleError(err)
        defer listen.Close()
        conn, err := listen.Accept()
        handleError(err)

        for {
                netData, err := bufio.NewReader(conn).ReadString('\n')
                handleError(err)
                if strings.TrimSpace(string(netData)) == "STOP" {
                        fmt.Println("Exiting TCP server!")
                        return
                }
                
                fmt.Print("-> ", string(netData))
                myTime := time.Now().Format(time.RFC3339) + "\n"
                conn.Write([]byte(myTime))
        }
}

func handleError(err error) {
        if err != nil {
                fmt.Println(err)
                return
        }
}