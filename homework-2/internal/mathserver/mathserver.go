package mathserver

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
	enigma   string
	answer   string
)

func InitMathServer() {
	listener, err := net.Listen("tcp", "localhost:8000")
	fmt.Println("Start chat server at: localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	newEnigma()

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

	ch <- "Set your name: "
	scannerName := bufio.NewScanner(conn)
	scannerName.Scan()
	name := scannerName.Text()
	ch <- "You are " + name
	messages <- name + ": " + "has arrived"
	entering <- ch

	messages <- "Solve the given expression: " + enigma

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + ": " + input.Text()
		// if user won
		if input.Text() == answer {
			messages <- name + " WON!"
			newEnigma()
			messages <- "Solve the given expression: " + enigma
		}
		// if user give up
		if strings.ToLower(input.Text()) == "give up" {
			messages <- name + " LOST! Answer: " + answer
			newEnigma()
			messages <- "Solve the given expression: " + enigma
		}
	}
	leaving <- ch
	messages <- name + " has left"
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

// newEnigma func generated new enigma
func newEnigma() {
	sourceRand := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(sourceRand)

	operand1 := randomizer.Intn(100)
	operand2 := randomizer.Intn(100)
	randomOperator := randomizer.Intn(4)

	switch randomOperator {
	case 0:
		enigma = strconv.Itoa(operand1) + "+" + strconv.Itoa(operand2)
		answer = strconv.Itoa(operand1 + operand2)
	case 1:
		enigma = strconv.Itoa(operand1) + "-" + strconv.Itoa(operand2)
		answer = strconv.Itoa(operand1 - operand2)
	case 2:
		enigma = strconv.Itoa(operand1) + "*" + strconv.Itoa(operand2)
		answer = strconv.Itoa(operand1 * operand2)
	case 3:
		enigma = strconv.Itoa(operand1) + "/" + strconv.Itoa(operand2)
		divNum := float64(operand1) / float64(operand2)
		answer = strconv.FormatFloat(divNum, 'f', -1, 64)
	}
}
