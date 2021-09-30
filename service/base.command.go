package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
)

func Command(cmd, natsUri string, natsOptions []nats.Option, args map[string]string) {
	defer func() {
		if r := recover(); r != nil {
			if _, e := fmt.Printf("%v", r); e != nil {
				panic(e)
			}
		}
	}()

	natsConn, err := nats.Connect(natsUri, natsOptions...)
	if err != nil {
		panic(err)
	}
	c, err = uniform.ConnectorNats(d, natsConn)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	d.Page(-1, traceRate, true, AppName, nil, "", "", nil, func(p diary.IPage) {
		p.Info(cmd, nil)
		if err := c.Publish(p, local(command(cmd)), uniform.Request{
			Parameters: args,
		}); err != nil {
			panic(err)
		}
	})

	_ = c.Drain()
}
