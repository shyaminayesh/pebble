package logger

import (
	"fmt"
)

func Println(color string, info string, msg string) {

	// WHITE
	if color == "white" {
		fmt.Println("\033[36m" + info + "\033[0m" + msg)
	}

	// GREEN
	if color == "green" {
		fmt.Println("\033[32m" + info + "\033[0m" + msg)
	}

	if color == "red" {
		fmt.Println("\033[31m" + info + "\033[0m" + msg)
	}

}