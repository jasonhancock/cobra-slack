package notifications

import (
	"fmt"
	"os"

	"github.com/jasonhancock/go-env"
	"github.com/jasonhancock/slack-go-webhook"
	"github.com/spf13/cobra"
)

const (
	envSlackWebhookURL = "SLACK_WEBHOOK_URL"
	envSlackChannel    = "SLACK_CHANNEL"
	envSlackBotName    = "SLACK_BOT_NAME"
)

// Config contains the configuration for Slack.
type Config struct {
	slackWebhookURL string
	slackChannel    string
	botName         string
}

// Enabled returns true if the required configuration parameters have been provided for sending webhooks to Slack.
func (o Config) Enabled() bool {
	return o.slackWebhookURL != "" && o.slackChannel != ""
}

// AddFlags adds the flags for slack to the command.
func AddFlags(cmd *cobra.Command, defaultBotName string) *Config {
	var opts Config

	cmd.Flags().StringVar(
		&opts.slackWebhookURL,
		"slack-webhook-url",
		os.Getenv(envSlackWebhookURL),
		"The Slack webhook URL. Can be set via "+envSlackWebhookURL+" environment variable.",
	)

	cmd.Flags().StringVar(
		&opts.slackChannel,
		"slack-channel",
		os.Getenv(envSlackChannel),
		"The Slack channel to send to. Example: #my-channel. Can be set via "+envSlackChannel+" environment variable.",
	)

	cmd.Flags().StringVar(
		&opts.botName,
		"slack-bot-name",
		env.String(envSlackBotName, defaultBotName),
		"The Slack bot's name. Can be set via "+envSlackBotName+" environment variable.",
	)

	return &opts
}

// Send sends the message.
func (o Config) Send(payload slack.Payload) error {
	if !o.Enabled() {
		return nil
	}
	payload.Username = o.botName
	payload.Channel = o.slackChannel

	// Slack package returns a slice of errors. Convert into a multierror
	if err := slack.Send(o.slackWebhookURL, "", payload); err != nil {
		return fmt.Errorf("sending slack notification: %w", err)
	}

	return nil
}
