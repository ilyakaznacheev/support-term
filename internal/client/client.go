package client

import (
	"strings"
	"time"

	// "log"
	"bufio"
	"fmt"
	"os"
	"os/signal"

	"github.com/ilyakaznacheev/support-term/internal/types"
	nats "github.com/nats-io/go-nats"
)

// Run start client support app
func Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your name")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		ec.Close()
		fmt.Println("Good luck, " + name)
	}()

	for {
		fmt.Println("Enter your question")

		question, _ := reader.ReadString('\n')
		msg := &types.Question{
			ID:       123, //RequestID(ec, "client-"+name),
			UserName: name,
			Text:     strings.TrimSpace(question),
		}

		resp := &types.Answer{}
		err := ec.Request("question", msg, resp, time.Minute)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", resp.SupName, resp.Text)
	}
}
