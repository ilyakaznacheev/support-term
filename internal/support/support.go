// Package support A support CLI termilal
package support

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/ilyakaznacheev/support-term/internal/types"
	nats "github.com/nats-io/go-nats"
)

type request struct {
	reply string
	msg   *types.Question
}

// Run start a CLI terminal for support users
func Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer func() {
		ec.Flush()
		ec.Close()
	}()

	questionCh := make(chan request)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your name")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("Work hard, " + name)

	ec.Subscribe("question", func(subject, reply string, msg *types.Question) {
		questionCh <- request{
			reply: reply,
			msg:   msg,
		}
	})
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-questionCh:
			fmt.Printf("%s: %s\nAnswer: ", msg.msg.UserName, msg.msg.Text)
			text, _ := reader.ReadString('\n')
			answer := &types.Answer{
				ID:      msg.msg.ID,
				SupName: name,
				Text:    strings.TrimSpace(text),
			}

			ec.Publish(msg.reply, answer)

		case <-interrupt:
			fmt.Println("Good job, " + name)
			return nil
		}
	}

}
