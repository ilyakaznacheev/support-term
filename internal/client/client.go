package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ilyakaznacheev/support-term/internal/pkg/message"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
)

type client struct {
	name      string
	ec        *nats.EncodedConn
	in        io.Reader
	out       io.Writer
	ctx       context.Context
	cancel    context.CancelFunc
	interrupt chan os.Signal
}

// NewClient create client app
func newClient(nc *nats.Conn, name string, in io.Reader, out io.Writer) *client {
	ec, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := &client{
		name:      name,
		ec:        ec,
		in:        in,
		out:       out,
		interrupt: interrupt,
	}

	return c
}

func (c *client) refreshContext() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go func() {
		<-c.interrupt
		c.cancel()
	}()
}

func (c *client) runMessageLoop() error {
	c.refreshContext()

	reader := bufio.NewReader(c.in)

	for {

		fmt.Fprintln(c.out, "Enter your question")

		question, _ := reader.ReadString('\n')
		msg := &message.Question{
			ID:       c.requestID(),
			UserName: c.name,
			Text:     strings.TrimSpace(question),
		}

		resp := &message.Answer{}
		err := c.ec.Request("question", msg, resp, time.Minute)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.out, "%s: %s\n", resp.SupName, resp.Text)

		if c.ctx.Err() != nil {
			break
		}
	}

	fmt.Fprintln(c.out, "Good luck, "+c.name)
	return nil
}

func (c *client) requestID() int64 {
	resp := &message.NextID{}
	err := c.ec.Request("id-request", nil, resp, time.Minute)
	if err != nil {
		return 0
	}

	return resp.GeneratedID
}

// Run start client support app
func Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	fmt.Println("Enter your name")
	name, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	cl := newClient(nc, strings.TrimSpace(name), os.Stdin, os.Stdout)

	return cl.runMessageLoop()
}
