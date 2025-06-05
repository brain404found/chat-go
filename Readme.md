# Chat Go

A real-time chat application built with gRPC and Go. This project demonstrates a simple but effective implementation of a chat service using gRPC streaming.

## Architecture

The project follows a client-server architecture using gRPC for communication:

### Components

- **gRPC Server**: Handles chat room creation, message broadcasting, and client connections
- **gRPC Client**: Demonstrates how to connect to the chat server and send/receive messages
- **Protocol Buffers**: Defines the service interface and message structures

### Key Features

- Real-time message broadcasting
- Multiple chat rooms support
- Concurrent client connections
- In-memory message storage
- Stream-based communication

## Project Structure

```
.
├── api/                    # Protocol Buffers definitions
│   └── chat_v1/
│       └── chat.proto      # Service and message definitions
├── cmd/                    # Application entry points
│   ├── grpc_client/        # Chat client implementation
│   └── grpc_server/        # Chat server implementation
├── internal/              # Private application code
│   └── service/           # Business logic implementation
├── pkg/                   # Public library code
│   └── chat_v1/           # Generated protobuf code
└── Makefile              # Build and development commands
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)
- Go plugins for protoc

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/chat-go.git
cd chat-go
```

2. Install dependencies:
```bash
go mod download
```

3. Generate protobuf code:
```bash
make generate
```

### Running the Application

1. Start the server:
```bash
make run-server
```

2. In a separate terminal, run the client:
```bash
make run-client
```

The client will automatically:
- Create a new chat room
- Connect two users (John and Jack)
- Send random messages every 5-7 seconds
- Display received messages in real-time

## API

### Service Methods

- `CreateChat`: Creates a new chat room
- `ConnectChat`: Establishes a streaming connection to a chat room
- `SendMessage`: Sends a message to a chat room

### Message Types

- `Message`: Contains the message content, sender, and timestamp
- `CreateChatResponse`: Returns the ID of the created chat room
- `ConnectChatRequest`: Contains chat room ID and username
- `SendMessageRequest`: Contains chat room ID and message content

## Development

### Available Make Commands

- `make generate`: Generate protobuf code
- `make run-server`: Run the gRPC server
- `make run-client`: Run the gRPC client
- `make lint`: Run linter
- `make test`: Run tests

### Code Generation

The project uses Protocol Buffers for API definition. After modifying the proto files, run:

```bash
make generate
```
