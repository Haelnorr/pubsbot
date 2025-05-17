package bot

import (
	"github.com/Haelnorr/pubsbot/pkg/config"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Contains the session objects for a created bot
type Bot struct {
	Session   *discordgo.Session
	Logger    *zerolog.Logger
	Config    *config.Config
	statusMsg string
}

// Create a new bot and start a session
func NewBot(
	l *zerolog.Logger,
	cfg *config.Config,
) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		return nil, errors.Wrap(err, "discordgo.New")
	}
	bot := &Bot{
		Session: session,
		Logger:  l,
		Config:  cfg,
	}
	return bot, nil
}
