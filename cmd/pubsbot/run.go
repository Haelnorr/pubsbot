package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sync"

	"github.com/Haelnorr/pubsbot/internal/discord/bot"
	"github.com/Haelnorr/pubsbot/internal/discord/startup"
	"github.com/Haelnorr/pubsbot/pkg/config"
	"github.com/Haelnorr/pubsbot/pkg/logging"

	"github.com/pkg/errors"
)

// Initializes and runs the server
func run(ctx context.Context, w io.Writer, args map[string]string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config, err := config.GetConfig(args)
	if err != nil {
		return errors.Wrap(err, "config.GetConfig")
	}

	// Setup the logfile
	var logfile *os.File = nil
	if config.LogOutput == "both" || config.LogOutput == "file" {
		logfile, err = logging.GetLogFile(config.LogDir)
		if err != nil {
			return errors.Wrap(err, "logging.GetLogFile")
		}
		defer logfile.Close()
	}

	// Setup the console writer
	var consoleWriter io.Writer
	if config.LogOutput == "both" || config.LogOutput == "console" {
		consoleWriter = w
	}

	// Setup the logger
	logger, err := logging.GetLogger(
		config.LogLevel,
		consoleWriter,
		logfile,
		config.LogDir,
	)
	if err != nil {
		return errors.Wrap(err, "logging.GetLogger")
	}
	logger.Info().Msg("Logger initialized")

	// Initialize the discord bot
	discordBot, err := bot.NewBot(
		logger,
		config,
	)
	if err != nil {
		return errors.Wrap(err, "bot.NewBot")
	}
	logger.Debug().Msg("Bot created")

	// Runs the discord bot
	go func() {
		logger.Info().Msg("Starting discord bot")
		if err := startup.Start(ctx, discordBot); err != nil {
			logger.Error().Err(err).Msg("Error running bot")
		}
	}()

	// Handles graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		defer cancel()
		if err := startup.Stop(discordBot); err != nil {
			logger.Error().Err(err).Msg("Error shutting down discord bot")
		}
	}()
	wg.Wait()
	logger.Info().Msg("Shutting down")
	return nil

}
