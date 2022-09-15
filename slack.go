package notifications

import (
	"fmt"
	"os"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/hashicorp/go-multierror"
	"github.com/jasonhancock/go-env"
	"github.com/spf13/cobra"
)

const (
	envSlackWebhookURL = "SLACK_WEBHOOK_URL"
	envSlackChannel    = "SLACK_CHANNEL"
	envSlackBotName    = "SLACK_BOT_NAME"
)

// SlackOptions contains the configuration for Slack.
type SlackOptions struct {
	slackWebhookURL string
	slackChannel    string
	botName         string
}

// Enabled returns true if the required configuration parameters have been provided for sending webhooks to Slack.
func (o SlackOptions) Enabled() bool {
	return o.slackWebhookURL != "" && o.slackChannel != ""
}

// AddSlackFlags adds the flags for slack to the command.
func AddSlackFlags(cmd *cobra.Command, defaultBotName string) *SlackOptions {
	var opts SlackOptions

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
func (o SlackOptions) Send(payload slack.Payload) error {
	if !o.Enabled() {
		return nil
	}
	payload.Username = o.botName
	payload.Channel = o.slackChannel

	// Slack package returns a slice of errors. Convert into a multierror
	var retErrs error
	if errs := slack.Send(o.slackWebhookURL, "", payload); len(errs) > 0 {
		for _, v := range errs {
			retErrs = multierror.Append(retErrs, v)
		}
	}
	if retErrs != nil {
		return fmt.Errorf("sending slack notification: %w", retErrs)
	}

	return nil
}
