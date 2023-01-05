package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	sql2 "github.com/ThreeDotsLabs/watermill-sql/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"postcard"
	"postcard/storage"
)

func main() {
	conn := "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := storage.NewDefaultSimplePostcardRepository(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}

	logger := watermill.NewStdLogger(false, false)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		log.Fatal(err)
	}

	sqlSub, err := sql2.NewSubscriber(db, sql2.SubscriberConfig{
		SchemaAdapter:    storage.WatermillSchemaAdapter{},
		OffsetsAdapter:   sql2.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
	}, logger)

	amqpURI := "amqp://guest:guest@localhost:5672/"
	amqpConfig := amqp.NewDurablePubSubConfig(amqpURI, func(topic string) string {
		return topic
	})

	amqpPub, err := amqp.NewPublisher(amqpConfig, logger)
	if err != nil {
		log.Fatal(err)
	}

	amqpSub, err := amqp.NewSubscriber(amqpConfig, logger)
	if err != nil {
		log.Fatal(err)
	}

	router.AddHandler(
		"sql-to-amqp",
		"events",
		sqlSub,
		"events",
		amqpPub,
		func(msg *message.Message) ([]*message.Message, error) {
			return []*message.Message{msg}, nil
		},
	)

	router.AddNoPublisherHandler(
		"read-amqp",
		"events",
		amqpSub,
		func(msg *message.Message) error {
			fmt.Printf("Received message %v: %v\n", msg.UUID, string(msg.Payload))
			return nil
		},
	)

	fwd, err := forwarder.NewForwarder(sqlSub, amqpPub, logger, forwarder.Config{
		ForwarderTopic: "events",
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := fwd.Run(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err = router.Run(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-fwd.Running()
	<-router.Running()

	pc, err := postcard.NewPostcard(uuid.NewString())
	if err != nil {
		log.Fatal(err)
	}

	err = repo.Save(context.Background(), pc)
	if err != nil {
		log.Fatal(err)
	}

	err = pc.Address(postcard.Address{
		Name:  "Someone",
		Line1: "Somewhere",
	}, postcard.Address{
		Name:  "Who",
		Line1: "Where",
	})
	if err != nil {
		log.Fatal(err)
	}

	err = repo.Save(context.Background(), pc)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
