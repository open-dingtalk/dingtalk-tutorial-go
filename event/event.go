package main

import (
	"context"
	"flag"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
	"time"
)

func OnEventReceived(_ context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
	eventHeader := event.NewEventHeaderFromDataFrame(df)

	logger.GetLogger().Infof("received event, delay=%s, eventType=%s, eventId=%s, eventBornTime=%d, eventCorpId=%s, eventUnifiedAppId=%s, data=%s",
		time.Duration(time.Now().UnixMilli()-eventHeader.EventBornTime)*time.Millisecond,
		eventHeader.EventType,
		eventHeader.EventId,
		eventHeader.EventBornTime,
		eventHeader.EventCorpId,
		eventHeader.EventUnifiedAppId,
		df.Data)

	response := payload.NewSuccessDataFrameResponse()
	if err := response.SetJson(event.NewEventProcessResultSuccess()); err != nil {
		return nil, err
	}
	return response, nil
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
	cli.RegisterAllEventRouter(OnEventReceived)

	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	select {}
}
