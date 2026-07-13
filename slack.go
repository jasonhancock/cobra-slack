package slack

import (
	"context"
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

	var conf Config

	cmd.Flags().StringVar(
		&conf.slackWebhookURL,
		"slack-webhook-url",
		os.Getenv(envSlackWebhookURL),
		"The Slack webhook URL. Can be set via "+envSlackWebhookURL+" environment variable.",
	)

	cmd.Flags().StringVar(
		&conf.slackChannel,
		"slack-channel",
		os.Getenv(envSlackChannel),
		"The Slack channel to send to. Example: #my-channel. Can be set via "+envSlackChannel+" environment variable.",
	)

	cmd.Flags().StringVar(
		&conf.botName,
		"slack-bot-name",
		env.String(envSlackBotName, defaultBotName),
		"The Slack bot's name. Can be set via "+envSlackBotName+" environment variable.",
	)

	return &conf
}

// Send sends the message.
func (c Config) Send(payload slack.Payload) error {
	if !c.Enabled() {
		return nil
	}
	payload.Username = c.botName
	payload.Channel = c.slackChannel

	// Slack package returns a slice of errors. Convert into a multierror
	if err := slack.Send(c.slackWebhookURL, "", payload); err != nil {
		return fmt.Errorf("sending slack notification: %w", err)
	}

	return nil
}

// Init initializes the config, only needs to be called if you're using a WebhookFetcher.
func (c *Config) Init(ctx context.Context, opts ...InitOption) error {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	if o.webhookFetcher == nil {
		return nil
	}

	// if the webhook url has been set or provided already, don't call the fetcher.
	if c.slackWebhookURL != "" {
		return nil
	}

	u, err := o.webhookFetcher(ctx)
	if err != nil {
		return err
	}

	c.slackWebhookURL = u
	return nil
}

type options struct {
	webhookFetcher WebhookFetcher
}

// InitOption is used to customize the initialization of the config.
type InitOption func(*options)

// WithWebhookFetcher sets a WebhookFetcher to use if the webhook URL isn't directly provided.
func WithWebhookFetcher(fetcher WebhookFetcher) InitOption {
	return func(o *options) {
		o.webhookFetcher = fetcher
	}
}

// WebhookFetcher allows you to fetch the webhook url from something like a
// secret management system. It should return the webhook URL string, or an empty
// string and an error. If you are using a WebhookFetcher, you need to manually call
// Init so that your fetcher gets called and initializes the configuration.
type WebhookFetcher func(ctx context.Context) (string, error)
