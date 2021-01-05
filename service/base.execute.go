package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
)

func Execute(test bool, natsUri, environment, level string, rate int, handler diary.H, argsMap M) {
	lvl := diary.ConvertFromTextLevel(level)
	if diary.IsValidLevel(lvl) {
		panic(fmt.Sprintf("level must be one of the following values: %s", strings.Join(diary.TextLevels, ", ")))
	}
	testMode = test
	traceRate = rate
	env = environment

	args = M{}
	if argsMap != nil {
		args = argsMap
	}
	args["nats"] = natsUri
	args["env"] = environment

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

	d = diary.Dear(AppClient, AppProject, AppName, nil, "git@github.com:go-uniform/uniform.git", AppCommit, nil, nil, lvl, handler)
	d.Page(-1, traceRate, true, AppName, nil, "", "", nil, func(p diary.IPage) {
		// subscribe all actions
		for topic, handler := range actions {
			p.Info("subscribe", diary.M{
				"project": AppProject,
				"topic":   topic,
				"handler": runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			})
			subscription, err := c.QueueSubscribe(fmt.Sprintf("%s.%s", AppProject, topic), AppName, handler)
			if err != nil {
				p.Error("subscribe", "failed to subscribe for topic", diary.M{
					"project": AppProject,
					"topic": topic,
					"error": err,
				})
			}
			subscriptions[topic] = subscription
		}

		Run(p)

		// Go signal notification works by sending `os.Signal`
		// values on a channel. We'll create a channel to
		// receive these notifications (we'll also make one to
		// notify us when the program can exit).
		signals := make(chan os.Signal, 1)

		// `signal.Notify` registers the given channel to
		// receive notifications of the specified signals.
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

		// The program will wait here until it gets the
		// expected signal (as indicated by the goroutine
		// above sending a value on `done`) and then exit.
		p.Notice("signal.wait", diary.M{
			"signals": []string{
				"syscall.SIGINT",
				"syscall.SIGTERM",
				"syscall.SIGKILL",
			},
		})
		sig := <-signals
		p.Notice("signal.received", diary.M{
			"signal": sig,
		})

		p.Notice("unsubscribe.all", diary.M{
			"topics.actions": reflect.ValueOf(actions).MapKeys(),
			"topics.subscriptions": reflect.ValueOf(subscriptions).MapKeys(),
			"count.actions": len(actions),
			"count.subscriptions": len(subscriptions),
		})

		// unsubscribe all actions
		for topic, subscription := range subscriptions {
			p.Notice("unsubscribe", diary.M{
				"topic": topic,
			})
			if err := subscription.Unsubscribe(); err != nil {
				p.Error("unsubscribe", "failed to unsubscribe from topic", diary.M{
					"topic": topic,
					"error": err,
				})
			}
		}

		p.Notice("drain", nil)
		// Drain connection (Preferred for responders)
		// Close() not needed if this is called.
		if err := c.Drain(); err != nil {
			// this error might not reach the diary.write topic listener since we are busy shutting down service
			// do not expect to see this message in the diary logs
			p.Error("drain", "failed to drain connection", diary.M{
				"error": err,
			})
		}
	})
}
