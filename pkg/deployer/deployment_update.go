package deployer

import (
	"context"
	"time"

	"github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

func (c *DefaultClient) updateDeployments(ctx context.Context) {

	inCh := make(chan Event)
	outCh := make(chan Event)

	go coalesce(inCh, outCh)

	go func() {
		bch := make(chan int, 1)
		bLast := 0
		for {
			c.statusCache.Register(bch, bLast)
			select {
			case <-ctx.Done():
				return
			case bLast = <-bch:
				e := Event(bLast)

				inCh <- e
			}
		}
	}()

	go func() {
		for event := range outCh {
			c.logger.Infow("status changed detected, updating deployments",
				"event", event,
			)

			err := c.updateDeploymentsStatus(ctx)
			if err != nil {
				c.logger.Errorw("failed to update deployment status",
					"error", err,
				)
			}
		}
	}()

	go func() {
		<-ctx.Done()
		close(inCh)
		close(outCh)
	}()
}

func (c *DefaultClient) updateDeploymentsStatus(rootCtx context.Context) error {

	var (
		err    error
		ctx    context.Context
		cancel context.CancelFunc
	)

	for k, v := range c.statusCache.List() {
		ctx, cancel = context.WithTimeout(rootCtx, time.Second*5)
		c.logger.Infow("updating deployment",
			"id", k,
			"status", v.Status(),
			"replicas", v.AvailableReplicas,
		)
		_, err = c.client.UpdateDeployment(ctx, &deployer_v1.UpdateDeploymentRequest{
			Id:                k,
			Status:            v.Status(),
			AvailableReplicas: v.AvailableReplicas,
		})
		if err != nil {
			c.logger.Errorw("failed to update deployment record in Dotscience",
				"error", err,
				"deployment_id", k,
			)
		}
		cancel()
	}

	return nil
}

type Event int

func NewEvent() Event                   { return 0 }
func (e Event) Merge(other Event) Event { return e + other }

func coalesce(in <-chan Event, out chan<- Event) {
	event := NewEvent()
	timer := time.NewTimer(0)

	var timerCh <-chan time.Time
	var outCh chan<- Event

	for {
		select {
		case e := <-in:
			event = event.Merge(e)
			if timerCh == nil {
				timer.Reset(1500 * time.Millisecond)
				timerCh = timer.C
			}
		case <-timerCh:
			outCh = out
			timerCh = nil
		case outCh <- event:
			event = NewEvent()
			outCh = nil
		}
	}
}
