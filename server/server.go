package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	FullName  string   `json:"fullname"`
	Emails    []string `json:"emails"`
	Addresses []string `json:"addresses"`
}

var users []User
var clientKeys map[*net.Conn]string

func loadUsers(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return err
	}

	return nil
}

func saveUsers(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		return err
	}

	return nil
}
func handleConnection(c net.Conn) {
	defer c.Close()

	var key string
	// Authentication
	authenticated := false
	for authenticated == false {

		username, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		username = strings.TrimSpace(username)

		password, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		password = strings.TrimSpace(password)

		for _, u := range users {
			if u.Username == username && u.Password == password {
				authenticated = true
				break
			}
		}

		if authenticated == false {
			fmt.Fprintf(c, "Invalid credentials. Please try again.\n")
		} else {

			// Generate a unique key for the authenticated user
			rand.Seed(time.Now().UnixNano())
			key = strconv.Itoa(rand.Intn(1000000))
			// Send successful signal
			fmt.Fprintf(c, "successful\n")
			// Send the key to the client
			clientRes, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(clientRes)
			fmt.Fprintf(c, "%s\n", key)
			//clientKeys[&c] = key
			break
		}
	}

	for {

		rand.Seed(time.Now().UnixNano())
		target := rand.Intn(100) + 1
		fmt.Printf("Target number is: %d\n", target)

		for {
			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			guess, err := strconv.Atoi(strings.TrimSpace(netData))

			if guess == -1 {
				fmt.Fprintf(c, "%s_GameOver\n", key)
				os.Exit(0)
			} else if guess < target {
				fmt.Fprintf(c, "%s_To Low\n", key)
			} else if guess > target {
				fmt.Fprintf(c, "%s_To Hight\n", key)
			} else {
				fmt.Fprintf(c, "%s_Correct\n", key)
				break // Break the inner loop to generate a new target number
			}
		}
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	// Load users from JSON file
	if err := loadUsers("users.json"); err != nil {
		fmt.Println("Error loading users:", err)
		return
	}

	// Initialize client keys map
	clientKeys = make(map[*net.Conn]string)

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
