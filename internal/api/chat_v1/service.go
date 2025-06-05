package chatv1

import (
	"chat-go/pkg/chat_v1"
	"sync"
)

type Chat struct {
	streams map[string]chat_v1.ChatV1_ConnectChatServer
	mu      sync.RWMutex
}

type Service struct {
	chat_v1.UnimplementedChatV1Server

	chats  map[string]*Chat
	muChat sync.RWMutex

	chMessage map[string]chan *chat_v1.Message
	muChannel sync.RWMutex
}

func NewService() *Service {
	return &Service{
		chats:     make(map[string]*Chat),
		chMessage: make(map[string]chan *chat_v1.Message),
	}
}
