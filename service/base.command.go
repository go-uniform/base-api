package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
)


func Command(command, natsUri string) {
	defer func() {
		if r := recover(); r != nil {
			if _, e := fmt.Printf("%v", r); e != nil {
				panic(e)
			}
		}
	}()

	natsConn, err := nats.Connect(natsUri)
	if err != nil {
		panic(err)
	}
	c, err = uniform.ConnectorNats(d, natsConn)
	if err != nil {
		panic(err)
	}

	// Close connection
	defer c.Close()

	d = diary.Dear(AppClient, AppProject, AppName, nil, "git@github.com:go-uniform/uniform.git", AppCommit, nil, nil, diary.LevelFatal, nil)
	d.Page(-1, traceRate, true, AppName, nil, "", "", nil, func(p diary.IPage) {
		fmt.Printf("executing %s command\n", command)
		if err := c.Publish(p, fmt.Sprintf("command.%s", command), uniform.Request{}); err != nil {
			panic(err)
		}
	})

	// Drain connection
	_ = c.Drain()
}
