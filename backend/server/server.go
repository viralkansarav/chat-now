package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/viralkansarav/chat-now/proto/chatpb"
	"google.golang.org/grpc"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[string]pb.ChatService_ChatStreamServer
}

func (s *chatServer) ChatStream(stream pb.ChatService_ChatStreamServer) error {
	var user string

	// so, first thing we do is receive the initial message to get the sender's name
	msg, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("oops, failed to receive initial message: %v", err)
	}
	user = msg.Sender
	// now add this user to our map of active clients
	s.mu.Lock()
	s.clients[user] = stream
	s.mu.Unlock()

	// we start a goroutine to handle incoming messages from this user
	go func() {
		for {
			// receiving a message from the client
			incomingMsg, err := stream.Recv()
			if err != nil {
				log.Printf("whoops, error receiving message: %v", err)
				break
			}

			log.Printf("received message from %s: %s", user, incomingMsg.Message)

			// once we get the message, we broadcast it to everyone connected
			s.broadcastMessage(incomingMsg)
		}
	}()

	// now we just wait for the connection to be done or closed
	for {
		select {
		case <-stream.Context().Done():
			// client closed the connection, so let's clean up
			log.Printf("connection with %s closed", user)
			s.mu.Lock()
			delete(s.clients, user)
			s.mu.Unlock()
			return nil
		}
	}
}

func (s *chatServer) broadcastMessage(msg *pb.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// we're gonna loop through all the connected clients and send them the message
	for _, client := range s.clients {
		if err := client.Send(msg); err != nil {
			log.Printf("error sending message to client: %v", err)
		}
	}
}

func main() {
	// Let's open the port 50051 for the server to listen
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("couldn't start listening: %v", err)
	}

	// now we create our gRPC server and register the ChatService
	s := grpc.NewServer()
	chatSrv := &chatServer{clients: make(map[string]pb.ChatService_ChatStreamServer)}

	// Register the ChatService server with the gRPC server
	pb.RegisterChatServiceServer(s, chatSrv)

	fmt.Println("server's running on port 50051...")
	// we now start serving, and it’ll keep going until something breaks
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
