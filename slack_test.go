package notifications

import (
	"testing"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/stretchr/testify/require"
)

func TestSlackUnconfigured(t *testing.T) {
	opts := SlackOptions{}
	err := opts.Send(slack.Payload{Text: "testing"})
	require.NoError(t, err)
}
