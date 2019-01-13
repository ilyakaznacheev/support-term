package client

import (
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

	questionCh := make(chan *types.Question, 1)
	ec.BindSendChan("question", questionCh)

	answerCh := make(chan *types.Answer)
	ec.BindRecvChan("answer", answerCh)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your name")
	name, _ := reader.ReadString('\n')

	for {
		fmt.Println("Enter your question")

		question, _ := reader.ReadString('\n')
		msg := types.Question{
			ID:       123,
			UserName: name,
			Text:     question,
		}

		questionCh <- &msg

		select {
		case answer := <-answerCh:
			fmt.Printf("%s: %s\n", answer.SupName, answer.Text)

		case <-interrupt:
			fmt.Println("Good luck, " + name)
			return nil
		}
	}
}
