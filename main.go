package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer server.Close()
	outC := make(chan string)
	conns := new([]net.Conn)
	go outputChannel(outC, conns)
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		} else {
			fmt.Println(conn.LocalAddr().String(), " has connected.")
			*conns = append(*conns, conn)
			go processClient(conn, outC)
		}
	}
}

func outputChannel(outC chan string, conns *[]net.Conn) {
	for {
		message := <-outC
		for i := 0; i < len(*conns); i++ {
			(*conns)[i].Write([]byte(message))
		}
	}
}

func processClient(connection net.Conn, outC chan string) {
	connection.Write([]byte("Enter your name: "))
	nameBuffr := make([]byte, 1024)
	nameLen, err := connection.Read(nameBuffr)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	name := string(nameBuffr[:nameLen-1])
	outC <- fmt.Sprintf("%s, has joined the chat.\n", name)

	for {
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if (err != nil && err == io.EOF) || string(buffer[:mLen-1]) == "exit" {
			connection.Close()
		  outC <- fmt.Sprintf("%s, has left the chat.\n", name)
			return
		}
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		outC <- fmt.Sprintf("%s: %s", name, string(buffer[:mLen]))
	}
}
