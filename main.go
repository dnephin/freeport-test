package freeporttest

import (
	"context"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/sdk/freeport"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

type TestingT interface {
	assert.TestingT
	Logf(format string, args ...interface{})
}

func RunTestConflicts(t TestingT, pkgName string, batchSize int) {
	g, ctx := errgroup.WithContext(context.Background())

	ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	var holdDuration = int64(200 * time.Millisecond)
	var loopDuration = int64(5 * time.Millisecond)
	var j = 2

	seed := time.Now().UnixNano()
	t.Logf("seed %v", seed)
	rnd := rand.New(rand.NewSource(seed))

	for i := 0; i < j; i++ {
		i := i
		g.Go(func() error {
			for ctx.Err() == nil {
				ports, err := freeport.Take(batchSize)
				if err != nil {
					return err
				}
				t.Logf("free %v %v: %v", pkgName, i, ports)

				listeners := make([]net.Listener, len(ports))
				for i, port := range ports {
					l, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
					if err != nil {
						return err
					}
					listeners[i] = l
				}
				time.Sleep(time.Duration(holdDuration + rnd.Int63n(holdDuration)))

				for _, l := range listeners {
					if err := l.Close(); err != nil {
						return err
					}
				}
				freeport.Return(ports)
				time.Sleep(time.Duration(loopDuration + rnd.Int63n(loopDuration)))
			}
			return nil
		})

		for b := 0; b < batchSize; b++ {
			g.Go(func() error {
				for ctx.Err() == nil {
					l, err := net.Listen("tcp", "127.0.0.1:0")
					if err != nil {
						return err
					}

					_, port, err := net.SplitHostPort(l.Addr().String())
					if err != nil {
						return err
					}
					t.Logf("zero %v %v: %v", pkgName, i, port)
					time.Sleep(time.Duration(holdDuration + rnd.Int63n(holdDuration)))

					if err := l.Close(); err != nil {
						return err
					}
					time.Sleep(time.Duration(loopDuration + rnd.Int63n(loopDuration)))
				}
				return nil
			})
		}
	}

	assert.NilError(t, g.Wait())
}
