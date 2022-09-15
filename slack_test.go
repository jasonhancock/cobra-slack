package notifications

import (
	"testing"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/stretchr/testify/require"
)

func TestUnconfigured(t *testing.T) {
	opts := Config{}
	err := opts.Send(slack.Payload{Text: "testing"})
	require.NoError(t, err)
}
