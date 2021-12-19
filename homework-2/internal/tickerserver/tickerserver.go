package tickerserver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	messages = make(chan string)
)

func InitTickerServer() {
	listener, err := net.Listen("tcp", "localhost:8000")
	fmt.Println("Start ticker server at: localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
		go messenger()
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	for {
		select {
		case msg := <-messages:
			_, err := io.WriteString(c, msg)
			if err != nil {
				return
			}
		default:
			_, err := io.WriteString(c, time.Now().Format("15:04:05\n\r"))
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// messenger func for send message to clients from server
func messenger() {
	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			return
		}
		messages <- line
	}
}
