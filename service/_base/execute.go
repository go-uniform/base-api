package _base

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"service/service/info"
	"strings"
	"sync"
	"syscall"
	"time"
)

func Execute(limit int, test, virtual bool, natsUri string, natsOptions []nats.Option, runBefore func(shutdown chan bool, group *sync.WaitGroup, p diary.IPage), runAfter func(shutdown chan bool, group *sync.WaitGroup, p diary.IPage), shutdownBefore func(shutdown chan bool, group *sync.WaitGroup, p diary.IPage), shutdownAfter func(shutdown chan bool, group *sync.WaitGroup, p diary.IPage)) {
	// set rate limiting duration using limit arg
	rateLimit := time.Nanosecond
	if limit > 0 && limit < 1000000 {
		rateLimit = time.Second / time.Duration(limit)
	}

	// set global TestMode flag based on test arg
	info.TestMode = test

	// set global Virtualize flag based on virtual arg
	info.Virtualize = virtual

	// connect to nats backbone
	natsConn, err := nats.Connect(natsUri, natsOptions...)
	if err != nil {
		panic(err)
	}
	info.Conn, err = uniform.ConnectorNats(info.Diary, natsConn)
	if err != nil {
		panic(err)
	}

	// on exit close the nats connection
	defer info.Conn.Close()

	info.Diary.Page(-1, info.TraceRate, true, info.AppName, nil, "", "", nil, func(p diary.IPage) {
		// a channel that will be closed when shutdown signal is received
		shutdown := make(chan bool)
		// a group used to check that all parallel running threads have been closed before shutdown routine starts
		group := &sync.WaitGroup{}

		p.Notice("startup", diary.M{
			"nats":       info.Args["nats"],
			"natsCert":   info.Args["natsCert"],
			"natsKey":    info.Args["natsKey"],
			"disableTls": info.Args["disableTls"],
			"lvl":        info.Args["lvl"],
			"rate":       info.Args["rate"],
			"limit":      info.Args["limit"],
			"test":       info.Args["test"],
			"virtual":    info.Args["virtual"],
		})

		// service custom run routine before subscribing actions
		p.Notice("run.before", nil)
		if err := p.Scope("run.before", func(p diary.IPage) {
			runBefore(shutdown, group, p)
		}); err != nil {
			panic(err)
		}

		// subscribe all actions [generic]
		p.Notice("subscribe", nil)
		for topic, handler := range actions {
			p.Info(fmt.Sprintf("subscribe.%s", topic), diary.M{
				"project": info.AppProject,
				"topic":   topic,
				"handler": runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			})
			subscription, err := info.Conn.QueueSubscribe(rateLimit, topic, info.AppService, handler)
			if err != nil {
				p.Error("subscribe", "failed to subscribe for topic", diary.M{
					"project": info.AppProject,
					"topic":   topic,
					"error":   err,
				})
			}
			subscriptions[topic] = subscription
		}

		// subscribe all actions [service specific]
		for topic, handler := range actions {
			if !strings.HasPrefix(topic, info.AppService+".") {
				// skip all non-routine topics
				continue
			}

			topic = fmt.Sprintf("%s.%s", info.AppName, topic)
			p.Info(fmt.Sprintf("subscribe.%s", topic), diary.M{
				"project": info.AppProject,
				"topic":   topic,
				"handler": runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			})
			subscription, err := info.Conn.QueueSubscribe(rateLimit, topic, info.AppService, handler)
			if err != nil {
				p.Error("subscribe", "failed to subscribe for topic", diary.M{
					"project": info.AppProject,
					"topic":   topic,
					"error":   err,
				})
			}
			subscriptions[topic] = subscription
		}

		// service custom run routine after subscribing actions
		p.Notice("run.after", nil)
		if err := p.Scope("run.after", func(p diary.IPage) {
			runAfter(shutdown, group, p)
		}); err != nil {
			panic(err)
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

		// service custom shutdown routine before trigger and unsubscribing actions
		p.Notice("shutdown.before", nil)
		if err := p.Scope("shutdown.before", func(p diary.IPage) {
			shutdownBefore(shutdown, group, p)
		}); err != nil {
			panic(err)
		}

		// trigger shutdown to notify all other threads
		p.Notice("shutdown", nil)
		close(shutdown)

		// give other threads 3 seconds to execute graceful shutdown otherwise just forcefully shutdown
		waitCh := make(chan bool)
		go func() {
			group.Wait()
			close(waitCh)
		}()
		select {
		case <-time.Tick(time.Second * 3):
			break
		case <-waitCh:
			break
		}

		p.Notice("unsubscribe.all", diary.M{
			"topics.actions":       reflect.ValueOf(actions).MapKeys(),
			"topics.subscriptions": reflect.ValueOf(subscriptions).MapKeys(),
			"count.actions":        len(actions),
			"count.subscriptions":  len(subscriptions),
		})

		// unsubscribe all actions
		for topic, subscription := range subscriptions {
			p.Info(fmt.Sprintf("unsubscribe.%s", topic), diary.M{
				"project": info.AppProject,
				"topic":   topic,
			})
			if err := subscription.Unsubscribe(); err != nil {
				p.Error("unsubscribe", "failed to unsubscribe from topic", diary.M{
					"topic": topic,
					"error": err,
				})
			}
		}

		// service custom shutdown routine after trigger and unsubscribing actions
		p.Notice("shutdown.after", nil)
		if err := p.Scope("shutdown.after", func(p diary.IPage) {
			shutdownAfter(shutdown, group, p)
		}); err != nil {
			panic(err)
		}

		p.Notice("drain", nil)
		_ = info.Conn.Drain()
	})
}
