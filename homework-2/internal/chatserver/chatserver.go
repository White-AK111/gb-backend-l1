package chatserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func InitChatServer() {
	listener, err := net.Listen("tcp", "localhost:8000")
	fmt.Println("Start chat server at: localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	ch := make(chan string)
	go clientWriter(conn, ch)

	// get username from client
	ch <- "Set your name: "
	scannerName := bufio.NewScanner(conn)
	scannerName.Scan()
	name := scannerName.Text()
	ch <- "You are " + name
	messages <- name + ": " + "has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + ": " + input.Text()
	}
	leaving <- ch
	messages <- name + " has left"
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
