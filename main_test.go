package freeporttest

import (
	"context"
	"math/rand"
	"net"

	"testing"
	"time"

	"github.com/hashicorp/consul/sdk/freeport"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestConflicts(t *testing.T) {
	g, ctx := errgroup.WithContext(context.Background())

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var freeportBatchSize int = 5
	var holdDuration = int64(200 * time.Millisecond)
	var loopDuration = int64(5 * time.Millisecond)
	var j = 3

	seed := time.Now().UnixNano()
	t.Logf("seed %v", seed)
	rnd := rand.New(rand.New(rand.NewSource(seed)))

	for i := 0; i < j; i++ {
		i := i
		g.Go(func() error {
			for ctx.Err() == nil {
				ports, err := freeport.Take(freeportBatchSize)
				assert.NilError(t, err)
				t.Logf("freeport %v: %v", i, ports)

				time.Sleep(time.Duration(holdDuration + rnd.Int63n(holdDuration)))
				freeport.Return(ports)
				time.Sleep(time.Duration(holdDuration + rnd.Int63n(loopDuration)))
			}
			return ctx.Err()
		})

		g.Go(func() error {
			for ctx.Err() == nil {
				for b := 0; b < freeportBatchSize; b++ {
					l, err := net.Listen("tcp", "127.0.0.1:0")
					assert.NilError(t, err)

					t.Logf("port0 %v: %v", i, l.Addr())
					time.Sleep(time.Duration(holdDuration + rnd.Int63n(holdDuration)))
					assert.NilError(t, l.Close())
				}
				time.Sleep(time.Duration(holdDuration + rnd.Int63n(loopDuration)))
			}
			return ctx.Err()
		})
	}

	assert.NilError(t, g.Wait())
}
