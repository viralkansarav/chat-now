# Bidirectional Streaming Chat with gRPC

Welcome to the Bidirectional Streaming Chat app! 🚀

This project demonstrates a real-time chat application using **gRPC** and **Go**. It uses bidirectional streaming to allow clients to send and receive messages in real time. So, when one user sends a message, all connected clients can receive it instantly – just like a real chat app! 😄

## Getting Started

### 1. **Clone the Repo**

First, you'll need to clone this repo to your local machine:

```bash
git clone https://github.com/viralkansarav/chat-now.git
cd chat-now
```
### 2. **Start the gRPC Server**

Now that everything is set up, start the gRPC server. It listens for connections on localhost:50051.

Run this command in your terminal:

```bash
go run server.go
```
You should see this message if the server is running correctly:

```bash
Server is running on port 50051...
```

### 3. **Start the Client**

To test this out with multiple clients, you need to start the client application in different terminal windows.

Run this command to start the first client:

```bash
go run client.go
```
The program will ask you for your name, and once you enter it, the client will join the chat room and send messages. 🎤

Now, open another terminal window, and start another client with the same command:

```bash
go run client.go
```

You can repeat this process for as many clients as you like. Each client will be able to send and receive messages in real-time. ✨

### 6. **Test Multiple Clients**
Send a Message: Once both clients are connected, you can start sending messages! Type a message into each terminal window. As soon as one client sends a message, the other client will instantly see it.
Close a Client: When you exit one of the client terminals by typing exit, the server will detect the disconnection, and that client will be removed from the chat.

### 7. **Shut Down the Server**
Once you're done testing, stop the server by pressing Ctrl + C in the terminal where it's running.

## How It Works
gRPC Server: We have a gRPC server running on port 50051, which listens for incoming chat messages. It uses bidirectional streaming to send and receive messages in real time. Each client is connected and gets messages broadcasted to them.

Client: Clients connect to the server using the gRPC client. They can send messages to the server, and the server broadcasts these messages to all connected clients. Each client runs in its own terminal, and they can send/receive messages with no delay.

Chat Messages: Every message sent is wrapped in a ChatMessage object, which contains the sender, message, and timestamp. This keeps things organized and allows for proper handling of messages.

