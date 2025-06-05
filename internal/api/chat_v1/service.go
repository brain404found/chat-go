package chatv1

import (
	chatpb "chat-go/pkg/chat_v1"
	"sync"
)

type Chat struct {
	streams map[string]chatpb.ChatV1_ConnectChatServer
	mu      sync.RWMutex
}

type Service struct {
	chatpb.UnimplementedChatV1Server

	chats  map[string]*Chat
	muChat sync.RWMutex

	chMessage map[string]chan *chatpb.Message
	muChannel sync.RWMutex
}

func NewService() *Service {
	return &Service{
		chats:     make(map[string]*Chat),
		chMessage: make(map[string]chan *chatpb.Message),
	}
}
