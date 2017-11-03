package slack

import (
	"github.com/nlopes/slack"
	"github.com/tylerconlee/slab/config"
)

var c = config.LoadConfig()

func Connect() {
	api := slack.New("YOUR_TOKEN_HERE")
}
