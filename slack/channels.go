package slack

import "github.com/tylerconlee/slab/datastore"

// Channel represents an individual Slack channel, used either for DMs or
// public usage, in which Slab has access to.
type Channel struct {
	ID string
}

var (
	// DMChannelList is a collection of the individual DM channels that Slab
	// has access to
	DMChannelList []Channel
	// ChannelList is a collectoin of the individual channels that Slab has
	// access to
	ChannelList []Channel
)

// GetChannel takes the event from RTM and determines if the channel is
// part of a DM with a user that just initiated Slab, or if it's in a Slab
// monitored channel.
func GetChannel(channel string) (chantype int) {
	for _, c := range DMChannelList {
		if channel == c.ID {
			return 1
		}
	}
	for _, c := range ChannelList {
		if channel == c.ID {
			return 2
		}
	}
	return 0
}

// AddChannel takes a channel and a channel type and adds it to the
// corresponding list.
// Types: 1 = DM Channel, 2 = Channel
func AddChannel(channel string, chantype int) {
	if chantype == 1 {
		DMChannelList = append(DMChannelList, Channel{ID: channel})
		datastore.SaveChannels(map[string]interface{}{
			"channel": channel,
			"private": true,
		})
	} else {
		ChannelList = append(ChannelList, Channel{ID: channel})
		datastore.SaveChannels(map[string]interface{}{
			"channel": channel,
			"private": false,
		})
	}
}

func LoadPublicChannels() {
	public, err := datastore.LoadChannels(false)
	if err != nil {
		log.Error("Error loading public channels", map[string]interface{}{
			"module": "slack",
			"error":  err,
		})
	}
	for _, p := range public {
		ChannelList = append(ChannelList, Channel{ID: p})
	}
}

func LoadPrivateChannels() {
	private, err := datastore.LoadChannels(true)
	if err != nil {
		log.Error("Error loading private channels", map[string]interface{}{
			"module": "slack",
			"error":  err,
		})
	}
	for _, p := range private {
		DMChannelList = append(DMChannelList, Channel{ID: p})
	}
}
