package chatv1

import (
	"chat-go/pkg/chat_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) ConnectChat(req *chat_v1.ConnectChatRequest, stream chat_v1.ChatV1_ConnectChatServer) error {
	s.muChannel.RLock()
	chatChan, ok := s.chMessage[req.GetChatId()]
	s.muChannel.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	s.muChat.Lock()
	if _, okChat := s.chats[req.GetChatId()]; !okChat {
		s.chats[req.GetChatId()] = &Chat{
			streams: make(map[string]chat_v1.ChatV1_ConnectChatServer),
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
