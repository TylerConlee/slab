package slack

type Channel struct {
	ID string
}

var (
	DMChannelList []Channel
	ChannelList   []Channel
)

// GetChannel takes the event from RTM and determines if the channel is
// part of a DM with a user that just initiated Slab, or if it's in a Slab
// monitored channel.
func getChannel(channel string) (chantype int) {
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

func AddChannel(channel string, chantype int) {
	if chantype == 1 {
		DMChannelList = append(DMChannelList, Channel{ID: channel})
	} else {
		ChannelList = append(ChannelList, Channel{ID: channel})
	}
}
