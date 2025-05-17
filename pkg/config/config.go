package config

import (
	"fmt"
	"os"

	"github.com/Haelnorr/pubsbot/pkg/logging"
	"github.com/Haelnorr/pubsbot/pkg/slapshotapi"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel          zerolog.Level              // Log level for global logging. Defaults to info
	LogOutput         string                     // "file", "console", or "both". Defaults to console
	LogDir            string                     // Path to create log files
	DiscordBotToken   string                     // Discord Bot Token
	SlapshotAPIConfig *slapshotapi.SlapAPIConfig // Config for the SlapshotAPI
}

// Load the application configuration and get a pointer to the Config object
func GetConfig(args map[string]string) (*Config, error) {
	godotenv.Load(".env")
	var (
		logLevel  zerolog.Level
		logOutput string
		valid     bool
	)

	if args["loglevel"] != "" {
		logLevel = logging.GetLogLevel(args["loglevel"])
	} else {
		logLevel = logging.GetLogLevel(GetEnvDefault("LOG_LEVEL", "info"))
	}
	if args["logoutput"] != "" {
		opts := map[string]string{
			"both":    "both",
			"file":    "file",
			"console": "console",
		}
		logOutput, valid = opts[args["logoutput"]]
		if !valid {
			logOutput = "console"
			fmt.Println(
				"Log output type was not parsed correctly. Defaulting to console only",
			)
		}
	} else {
		logOutput = GetEnvDefault("LOG_OUTPUT", "console")
	}
	if logOutput != "both" && logOutput != "console" && logOutput != "file" {
		logOutput = "console"
	}
	slapapikey := os.Getenv("SLAPSHOT_API_KEY")
	slapapicfg, err := slapshotapi.NewSlapAPIConfig(
		GetEnvDefault("SLAPSHOT_API_ENV", "staging"),
		slapapikey,
	)
	if err != nil {
		return nil, errors.Wrap(err, "slapshotapi.NewSlapAPIConfig")
	}

	config := &Config{
		LogLevel:          logLevel,
		LogOutput:         logOutput,
		LogDir:            GetEnvDefault("LOG_DIR", ""),
		DiscordBotToken:   os.Getenv("DISCORD_BOT_TOKEN"),
		SlapshotAPIConfig: slapapicfg,
	}

	if config.DiscordBotToken == "" && args["dbver"] != "true" {
		return nil, errors.New("Envar not set: DISCORD_BOT_TOKEN")
	}
	if slapapikey == "" && args["dbver"] != "true" {
		return nil, errors.New("Envar not set: SLAPSHOT_API_KEY")
	}

	return config, nil
}
