package service

import (
	chatpb "chat-go/pkg/chat_v1"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) ConnectChat(req *chatpb.ConnectChatRequest, stream chatpb.ChatV1_ConnectChatServer) error {
	s.muChannel.RLock()
	chatChan, ok := s.chMessage[req.GetChatId()]
	s.muChannel.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	s.muChat.Lock()
	if _, okChat := s.chats[req.GetChatId()]; !okChat {
		s.chats[req.GetChatId()] = &Chat{
			streams: make(map[string]chatpb.ChatV1_ConnectChatServer),
		}
	}
	s.muChat.Unlock()

	s.chats[req.GetChatId()].mu.Lock()
	s.chats[req.GetChatId()].streams[req.GetUsername()] = stream
	s.chats[req.GetChatId()].mu.Unlock()

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}

			for _, st := range s.chats[req.GetChatId()].streams {
				if err := st.Send(msg); err != nil {
					return err
				}
			}
		case <-stream.Context().Done():
			s.chats[req.GetChatId()].mu.Lock()
			delete(s.chats[req.GetChatId()].streams, req.GetUsername())
			s.chats[req.GetChatId()].mu.Unlock()
			return nil
		}
	}
}

func (s *Service) CreateChat(ctx context.Context, _ *emptypb.Empty) (*chatpb.CreateChatResponse, error) {
	chatID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	s.chMessage[chatID.String()] = make(chan *chatpb.Message, 100)

	return &chatpb.CreateChatResponse{
		ChatId: chatID.String(),
	}, nil
}

func (s *Service) SendMessage(ctx context.Context, req *chatpb.SendMessageRequest) (*emptypb.Empty, error) {
	s.muChannel.RLock()
	chatChan, ok := s.chMessage[req.GetChatId()]
	s.muChannel.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	chatChan <- req.GetMessage()

	return &emptypb.Empty{}, nil
}