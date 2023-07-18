package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
        CONNECT := HOST+":"+PORT
        conn, err := net.Dial(TYPE, CONNECT)
        handleError(err)

        for {
                reader := bufio.NewReader(os.Stdin)
                fmt.Print(">> ")
                text, _ := reader.ReadString('\n')
                fmt.Fprintf(conn, text+"\n")

                message, _ := bufio.NewReader(conn).ReadString('\n')
                fmt.Print("->: " + message)
                if strings.TrimSpace(string(text)) == "STOP" {
                        fmt.Println("TCP client exiting...")
                        return
                }
	}
}

func handleError(err error) {
        if err != nil {
                fmt.Println(err)
                return
        }
}