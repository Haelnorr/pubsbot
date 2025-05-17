package bot

import (
	"context"
	"fmt"
	"github.com/Haelnorr/pubsbot/pkg/slapshotapi"
	"time"

	"github.com/pkg/errors"
)

func (b *Bot) updateStatus(ctx context.Context) error {
	queue, err := b.getPubsQueue(ctx)
	if err != nil {
		return errors.Wrap(err, "b.getPubsQueue")
	}

	msg := "In Queue: %v | In Match: %v"
	msg = fmt.Sprintf(msg, queue.InQueue, queue.InMatch)
	if msg == b.statusMsg {
		b.Logger.Debug().Msg("Status message not changed, not updating")
		return nil
	}
	err = b.Session.UpdateCustomStatus(msg)
	if err != nil {
		return errors.Wrap(err, "b.Session.UpdateCustomStatus")
	}
	b.statusMsg = msg
	b.Logger.Debug().Msg("Status message updated")
	return nil
}

func (b *Bot) getPubsQueue(ctx context.Context) (*slapshotapi.PubsQueue, error) {
	regions := make([]string, 1)
	regions[0] = slapshotapi.RegionOCEEast
	queue, err := slapshotapi.GetQueueStatus(ctx, regions, b.Config.SlapshotAPIConfig)
	if err != nil {
		return nil, errors.Wrap(err, "slapshotapi.GetQueueStatus")
	}
	return queue, nil
}

func (b *Bot) StartWatchingQueue(ctx context.Context) {
	b.Logger.Info().Msg("Queue watch has been started.")
	ticker := time.NewTicker(15 * time.Second)
	stoppedByContext := false
	go func() {
		defer func() {
			ticker.Stop()
			if !stoppedByContext {
				b.Logger.Warn().Msg("Queue watch has been stopped. Did an error occur?")
			}
		}()
		for {
			select {
			case <-ctx.Done():
				stoppedByContext = true
				b.Logger.Info().Msg("Stopping queue watch due to shutdown.")
				return
			case <-ticker.C:
				err := b.updateStatus(ctx)
				if err != nil {
					b.Logger.Error().Err(err).Msg("Error occured updating the queue status")
				}
			}
		}
	}()
}
