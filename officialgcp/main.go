package main

import (
	"MRE_opencensus_pubsub/common"
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"go.opencensus.io/trace"
)

var gcpPs *pubsub.Client
var gcpTopic *pubsub.Topic

func initPubSub(config *common.Config) error {
	var ctx = context.Background()
	var err error

	if gcpPs, err = pubsub.NewClient(ctx, config.GCPProjectID); err != nil {
		return err
	}
	gcpTopic = gcpPs.Topic(config.Topic)
	gcpTopic.PublishSettings = pubsub.PublishSettings{
		CountThreshold: 1,
	}
	return nil
}

func main() {
	ctx := context.Background()

	config, err := common.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}
	if err := initPubSub(config); err != nil {
		log.Fatalln(err)
	}
	if err := common.InitTracing(config.GCPProjectID); err != nil {
		log.Fatalln(err)
	}

	if err := handleRequest(ctx); err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Second * 10)
}

func handleRequest(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "main()")
	defer span.End()

	doSomething(ctx)
	if err := publishMessage(ctx); err != nil {
		return err
	}
	return nil
}

func doSomething(ctx context.Context) {
	_, span := trace.StartSpan(ctx, "doSomething()")
	defer span.End()

	time.Sleep(time.Second * 2)
	log.Printf("Did something.\n")
}

func publishMessage(ctx context.Context) error {
	publishResult := gcpTopic.Publish(ctx, &pubsub.Message{Data: []byte{1}})
	id, err := publishResult.Get(ctx)
	if err != nil {
		return err
	}
	log.Printf("Published message %s.\n", id)
	return nil
}
