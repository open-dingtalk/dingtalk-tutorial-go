package main

import (
	"context"
	"flag"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"message_api/pkg/receiver"
)

func main() {
	var corpId, clientId, clientSecret string
	flag.StringVar(&corpId, "corp_id", "", "your-corp-id")
	flag.StringVar(&clientId, "client_id", "", "your-client-id")
	flag.StringVar(&clientSecret, "client_secret", "", "your-client-secret")
	flag.Parse()
	if len(clientId) == 0 || len(clientSecret) == 0 {
		panic("command line options --client_id and --client_secret required")
	}

	logger.SetLogger(logger.NewStdTestLogger())

	cli := client.NewStreamClient(
		client.WithAppCredential(client.NewAppCredentialConfig(clientId, clientSecret)),
		client.WithOpenApiHost("https://pre-api.dingtalk.com"),
	)

	agiPlugin := receiver.NewAgiPlugin(corpId, clientId, clientSecret)
	cli.RegisterPluginCallbackRouter(agiPlugin.OnIncomingRequest)

	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	select {}
}
