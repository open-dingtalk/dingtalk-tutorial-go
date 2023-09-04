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
	replyMsg := "echo received message:\n"
	for _, line := range strings.Split(data.Text.Content, "\n") {
		replyMsg += fmt.Sprintf("\n> 1. %s", strings.TrimSpace(line))
	}
	replier := chatbot.NewChatbotReplier()
	if err := replier.SimpleReplyMarkdown(ctx, data.SessionWebhook, []byte("stream-tutorial-go"), []byte(replyMsg)); err != nil {
		return nil, err
	}
	return []byte(""), nil
}

func main() {
	var clientId, clientSecret string
	flag.StringVar(&clientId, "client_id", "", "your-client-id")
	flag.StringVar(&clientSecret, "client_secret", "", "your-client-secret")
	flag.Parse()

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
