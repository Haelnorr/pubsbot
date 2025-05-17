package startup

import (
	"context"
	"github.com/Haelnorr/pubsbot/internal/discord/bot"
	"time"

	"github.com/pkg/errors"
)

// Start the bot
func Start(ctx context.Context, b *bot.Bot) error {
	starttime := time.Now()
	err := b.Session.Open()
	if err != nil {
		return errors.Wrap(err, "b.session.Open")
	}

	// Start the queue watching
	b.StartWatchingQueue(ctx)

	b.Logger.Info().Dur("startup_time", time.Since(starttime)).Msg("Bot startup complete!")
	return nil
}

// Stop the bot
func Stop(b *bot.Bot) error {
	err := b.Session.Close()
	if err != nil {
		return err
	}
	return nil
}
