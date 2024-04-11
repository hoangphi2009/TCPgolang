package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// const (
// 	HOST = "localhost"
// 	PORT = "8080"
// 	TYPE = "tcp"
// )

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	var key string
	for {
		// Authentication
		fmt.Print("Enter username: ")
		username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(c, username)

		fmt.Print("Enter password: ")
		password, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(c, password)

		// Read the response from the server
		response, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("Server connection: " + response)
		c.Write([]byte("Client Connected\n"))
		if strings.Contains(response, "successful") {
			key, err = bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("You are now connected to the server. Your key is %s \n", key)
			break
		}

	}

	for {

		fmt.Print("Guess the number: ")
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n') // Read input from the user
		fmt.Fprintf(c, text)                                  // Send the input to the server

		message, _ := bufio.NewReader(c).ReadString('\n') // Read the response from the server
		key = strings.TrimSpace(key)
		fmt.Print("Server: " + message)
		if strings.TrimSpace(message) == key+"_"+"Correct" {
			fmt.Println("Congratulations! You've guessed the number correctly.\nAutomatically play with other guessing number, (PRESS -1 to end game) ?")
		} else if strings.TrimSpace(message) == key+"_"+"GameOver" {
			fmt.Println("Game Over! You've quit the game.")
			break
		}
	}
}
