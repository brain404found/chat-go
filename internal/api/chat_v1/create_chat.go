package chatv1

import (
	"chat-go/pkg/chat_v1"
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CreateChat(ctx context.Context, _ *emptypb.Empty) (*chat_v1.CreateChatResponse, error) {
	chatID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	s.chMessage[chatID.String()] = make(chan *chat_v1.Message, 100)

	return &chat_v1.CreateChatResponse{
		ChatId: chatID.String(),
	}, nil
}
