package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"strings"
)

func OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	replyMsg := []byte(fmt.Sprintf("echo received message: [%s]", strings.TrimSpace(data.Text.Content)))

	replier := chatbot.NewChatbotReplier()
	if err := replier.SimpleReplyText(ctx, data.SessionWebhook, replyMsg); err != nil {
		return nil, err
	}
	return []byte(""), nil
}

func main() {
	var clientId, clientSecret string
	flag.StringVar(&clientId, "client_id", "", "your-client-id")
	flag.StringVar(&clientSecret, "client_secret", "", "your-client-secret")
	flag.Parse()
	if len(clientId) == 0 || len(clientSecret) == 0 {
		panic("command line options --client_id and --client_secret required")
	}

	logger.SetLogger(logger.NewStdTestLogger())

	cli := client.NewStreamClient(client.WithAppCredential(client.NewAppCredentialConfig(clientId, clientSecret)))
	cli.RegisterChatBotCallbackRouter(OnChatBotMessageReceived)

	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	select {}
}
