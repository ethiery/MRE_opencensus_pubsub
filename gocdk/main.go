package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"MRE_opencensus_pubsub/common"

	"go.opencensus.io/trace"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/gcppubsub"
)

var topic *pubsub.Topic

func initPubSub(ctx context.Context, config *common.Config) error {
	var err error
	urlstr := fmt.Sprintf("gcppubsub://%s/%s", config.GCPProjectID, config.Topic)
	topic, err = pubsub.OpenTopic(ctx, urlstr)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()

	config, err := common.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}
	if err := initPubSub(ctx, config); err != nil {
		log.Fatalln(err)
	}
	if err := common.InitTracing(config.GCPProjectID); err != nil {
		log.Fatalln(err)
	}

	if err := handleRequest(ctx); err != nil {
		log.Fatalln(err)
	}

	if err := topic.Shutdown(ctx); err != nil {
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
	err := topic.Send(ctx, &pubsub.Message{Body: []byte{1}})
	if err != nil {
		return err
	}
	log.Printf("Published message\n")
	return nil
}
