// Package support A support CLI termilal
package support

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	"github.com/ilyakaznacheev/support-term/internal/types"
	nats "github.com/nats-io/go-nats"
)

// Run start a CLI terminal for support users
func Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	questionCh := make(chan *types.Question)
	ec.BindRecvChan("question", questionCh)

	answerCh := make(chan *types.Answer, 1)
	ec.BindSendChan("answer", answerCh)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your name")
	name, _ := reader.ReadString('\n')
	fmt.Println("Work hard, " + name)

	for {
		select {
		case msg := <-questionCh:
			fmt.Printf("%s: %s\nAnswer: ", msg.UserName, msg.Text)
			text, _ := reader.ReadString('\n')
			answer := types.Answer{
				ID:      msg.ID,
				SupName: name,
				Text:    text,
			}
			answerCh <- &answer

		case <-interrupt:
			fmt.Println("Good job, " + name)
			return nil
		}
	}
}
