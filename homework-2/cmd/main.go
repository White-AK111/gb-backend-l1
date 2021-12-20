package main

import (
	"fmt"
	"os"

	"github.com/White-AK111/gb-backend-l1/homework-2/internal/chatserver"
	"github.com/White-AK111/gb-backend-l1/homework-2/internal/client"
	"github.com/White-AK111/gb-backend-l1/homework-2/internal/mathserver"
	"github.com/White-AK111/gb-backend-l1/homework-2/internal/tickerserver"
)

func main() {

	if len(os.Args) == 2 {
		switch arg := os.Args[1]; arg {
		case "client":
			client.InitClient()
		case "chatserver":
			chatserver.InitChatServer()
		case "tickerserver":
			tickerserver.InitTickerServer()
		case "mathserver":
			mathserver.InitMathServer()
		default:
			fmt.Println("Incorrect argument, use \"client\" or \"chatserver\" or \"tickerserver\" or \"mathserver\"")
		}
	} else {
		fmt.Println("Use argument \"client\" or \"chatserver\" or \"tickerserver\" or \"mathserver\"")
	}
}
