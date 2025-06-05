package main

import (
	chatpb "chat-go/pkg/chat_v1"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := chatpb.NewChatV1Client(conn)

	chatID, err := createChat(ctx, client)
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}

	log.Printf(fmt.Sprintf("%s: %s\n", color.GreenString("Chat created"), color.YellowString(chatID)))

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Подключаемся к чату от имени пользователя John
	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, "John", 5*time.Second)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	// Подключаемся к чату от имени пользователя Jack
	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, "Jack", 7*time.Second)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	wg.Wait()
}

func connectChat(ctx context.Context, client chatpb.ChatV1Client, chatID string, username string, period time.Duration) error {
	stream, err := client.ConnectChat(ctx, &chatpb.ConnectChatRequest{
		ChatId:   chatID,
		Username: username,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[%v] - [from: %s]: %s\n",
				color.YellowString(message.GetCreatedAt().AsTime().Format(time.RFC3339)),
				color.BlueString(message.GetFrom()),
				message.GetText(),
			)
		}
	}()

	for {
		/*
			Ниже пример того, как можно считывать сообщения из консоли
			в демонстрационных целях будем засылать в чат рандомный текст раз в 5 секунд
			scanner := bufio.NewScanner(os.Stdin)
			var lines strings.Builder

			for {
				scanner.Scan()
				line := scanner.Text()
				if len(line) == 0 {
					break
				}

				lines.WriteString(line)
				lines.WriteString("\n")
			}

			err = scanner.Err()
			if err != nil {
				log.Println("failed to scan message: ", err)
			}
		*/

		time.Sleep(period)

		text := gofakeit.Word()

		_, err = client.SendMessage(ctx, &chatpb.SendMessageRequest{
			ChatId: chatID,
			Message: &chatpb.Message{
				From:      username,
				Text:      text,
				CreatedAt: timestamppb.Now(),
			},
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}

func createChat(ctx context.Context, client chatpb.ChatV1Client) (string, error) {
	res, err := client.CreateChat(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return res.GetChatId(), nil
}
