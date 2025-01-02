package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    startTCPClient()
}

func startTCPClient() {
    // Connect to the TCP server
    conn, err := net.Dial("tcp", "localhost:8000")
    if err != nil {
        fmt.Println("Failed to connect to server:", err)
        return
    }
    defer conn.Close()

    fmt.Println("Connected to the TCP server.")
    // Get the user's username and password
    username, password := getUsernameAndPasswordFromInput()

    // Send the username and password to the server
    _, err = conn.Write([]byte(fmt.Sprintf("%s %s join\n", username, password)))
    if err != nil {
        fmt.Println("Error writing to connection:", err)
        return
    }

    // Start a goroutine to read responses from the server
    go readResponsesFromServer(conn)

    // Read user input and send it to the server
    readAndSendMoves(conn, username, password)
}

func getUsernameAndPasswordFromInput() (string, string) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter your username: ")
    username, _ := reader.ReadString('\n')
    fmt.Print("Enter your password: ")
    password, _ := reader.ReadString('\n')
    return strings.TrimSpace(username), strings.TrimSpace(password)
}

func readAndSendMoves(conn net.Conn, username, password string) {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter move direction (up, down, left, right): ")
        text, _ := reader.ReadString('\n')
        if strings.TrimSpace(text) == "exit" {
            break
        }

        // Append the username, password, and move direction to the input text
        text = fmt.Sprintf("%s %s move %s", username, password, strings.TrimSpace(text))

        // Send the message to the server
        _, err := conn.Write([]byte(text + "\n"))
        if err != nil {
            fmt.Println("Error writing to connection:", err)
            return
        }
    }
}

func readResponsesFromServer(conn net.Conn) {
    buf := make([]byte, 1024)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading from connection:", err)
            return
        }
        response := strings.TrimSpace(string(buf[:n]))
        fmt.Println(response)
        if strings.Contains(response, "Login successful") {
            fmt.Println("You have successfully logged in. You can now join the game.")
        }
    }
}