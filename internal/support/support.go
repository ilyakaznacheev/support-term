// Package support A support CLI termilal
package support

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/ilyakaznacheev/support-term/internal/pkg/message"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
)

type request struct {
	reply string
	msg   *message.Question
}

type support struct {
	name      string
	ec        *nats.EncodedConn
	in        io.Reader
	out       io.Writer
	ctx       context.Context
	cancel    context.CancelFunc
	interrupt chan os.Signal
}

func newSupport(nc *nats.Conn, name string, in io.Reader, out io.Writer) *support {
	ec, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	s := &support{
		name:      name,
		ec:        ec,
		in:        in,
		out:       out,
		interrupt: interrupt,
	}

	return s
}

func (s *support) runMessageLoop() error {
	questionCh := make(chan request)

	reader := bufio.NewReader(s.in)

	s.ec.Subscribe("question", func(subject, reply string, msg *message.Question) {
		questionCh <- request{
			reply: reply,
			msg:   msg,
		}
	})

	for {
		select {
		case msg := <-questionCh:
			fmt.Fprintf(s.out, "%s: %s\nAnswer: ", msg.msg.UserName, msg.msg.Text)
			text, _ := reader.ReadString('\n')

			answer := &message.Answer{
				ID:      msg.msg.ID,
				SupName: s.name,
				Text:    strings.TrimSpace(text),
			}

			s.ec.Publish(msg.reply, answer)

		case <-s.interrupt:
			fmt.Fprintf(s.out, "Good job, "+s.name)
			s.ec.Flush()
			return nil
		}
	}
}

// Run start a CLI terminal for support users
func Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer func() {
		nc.Close()
	}()

	fmt.Println("Enter your name")
	name, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("Work hard, " + name)

	sp := newSupport(nc, strings.TrimSpace(name), os.Stdin, os.Stdout)

	return sp.runMessageLoop()
}
