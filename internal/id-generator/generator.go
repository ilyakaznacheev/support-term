// Package generator ID generation microservice
package generator

import (
	"log"
	"os"
	"os/signal"

	"github.com/ilyakaznacheev/support-term/internal/pkg/message"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
)

var currentID int64

// Run start ID genarator
func Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	ec, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	defer ec.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("Starting ID Generator")

	ec.Subscribe("id-request", func(subject, reply string, msg interface{}) {
		newID := &message.NextID{
			GeneratedID: nextID(),
		}
		log.Printf("Generated id: %d\n", newID.GeneratedID)
		ec.Publish(reply, newID)
	})
	if err != nil {
		return err
	}

	<-interrupt
	log.Println("Exiting ID Generator")
	return nil
}

func nextID() int64 {
	currentID++
	return currentID
}
