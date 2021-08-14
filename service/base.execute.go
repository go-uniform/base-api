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
	"time"
)

func Execute(limit int, test bool, natsUri string, natsOptions []nats.Option, argsMap M) {
	rateLimit := time.Nanosecond
	if limit > 0 && limit < 1000000 {
		rateLimit = time.Second / time.Duration(limit)
	}

	testMode = test

	args = M{}
	if argsMap != nil {
		args = argsMap
	}
	args["nats"] = natsUri

	natsConn, err := nats.Connect(natsUri, natsOptions...)
	if err != nil {
		panic(err)
	}
	c, err = uniform.ConnectorNats(d, natsConn)
	if err != nil {
		panic(err)
	}

	// Close connection
	defer c.Close()

	d.Page(-1, traceRate, true, AppName, nil, "", "", nil, func(p diary.IPage) {
		// run service custom startup routine before subscribing any actions
		if err := p.Scope("run", func(p diary.IPage) {
			Run(p)
		}); err != nil {
			panic(err)
		}

		// subscribe all actions [generic]
		for topic, handler := range actions {
			p.Info(fmt.Sprintf("subscribe.%s", topic), diary.M{
				"project": AppProject,
				"topic":   topic,
				"handler": runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			})
			subscription, err := c.QueueSubscribe(rateLimit, topic, AppService, handler)
			if err != nil {
				p.Error("subscribe", "failed to subscribe for topic", diary.M{
					"project": AppProject,
					"topic": topic,
					"error": err,
				})
			}
			subscriptions[topic] = subscription
		}

		// subscribe all actions [service specific]
		for topic, handler := range actions {
			if !strings.HasPrefix(topic, AppService + ".") {
				// skip all non-routine topics
				continue
			}

			topic = fmt.Sprintf("%s.%s", AppName, topic)
			p.Info(fmt.Sprintf("subscribe.%s", topic), diary.M{
				"project": AppProject,
				"topic":   topic,
				"handler": runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			})
			subscription, err := c.QueueSubscribe(rateLimit, topic, AppService, handler)
			if err != nil {
				p.Error("subscribe", "failed to subscribe for topic", diary.M{
					"project": AppProject,
					"topic": topic,
					"error": err,
				})
			}
			subscriptions[topic] = subscription
		}

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
