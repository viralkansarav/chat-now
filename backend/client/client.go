package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	pb "github.com/viralkansarav/chat-now/proto/chatpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// First, we try connecting to the server on localhost at port 50051
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// if we can't connect, we just print an error
		fmt.Printf("failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	// Start the bidirectional stream, we’re about to send and receive messages
	stream, err := client.ChatStream(context.Background())
	if err != nil {
		// something went wrong starting the stream, so we’ll let the user know
		fmt.Printf("failed to start stream: %v\n", err)
		return
	}

	// we need to read user input, so we set up a scanner to grab what they type
	scanner := bufio.NewScanner(os.Stdin)

	// let’s ask for their name first
	fmt.Print("Enter your name: ")
	var username string
	fmt.Scanln(&username)

	// send a "join chat" message to let others know you’re here
	initialMessage := &pb.ChatMessage{
		Sender:    username,
		Message:   "has joined the chat!",
		Timestamp: time.Now().String(),
	}
	// if sending this message fails, we gotta let the user know and exit
	if err := stream.Send(initialMessage); err != nil {
		fmt.Printf("failed to send initial message: %v\n", err)
		return
	}

	// now we start a goroutine that’ll listen for incoming messages from other clients
	go func() {
		for {
			// waiting for messages to come in from the server
			msg, err := stream.Recv()
			if err != nil {
				// if there’s an error receiving the message, we print it and stop
				fmt.Printf("error receiving message: %v\n", err)
				return
			}
			// print the received message in a nice format
			fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Sender, msg.Message)
		}
	}()

	// now let’s handle sending messages typed by the user
	for {
		// we prompt the user to type their message
		fmt.Print("You: ")
		scanner.Scan()
		message := scanner.Text()

		// if the user types 'exit', we break out of the loop and quit
		if strings.ToLower(message) == "exit" {
			break
		}

		// create a new chat message with the user's input
		chatMessage := &pb.ChatMessage{
			Sender:    username,
			Message:   message,
			Timestamp: time.Now().String(),
		}
		// send the message to the server, and check if it works
		if err := stream.Send(chatMessage); err != nil {
			// if sending fails, we print the error and exit
			fmt.Printf("failed to send message: %v\n", err)
			break
		}
	}
}
