package slack

import (
	"testing"
)

func TestGetChannel(t *testing.T) {
	type args struct {
		channel string
	}
	DMChannelList = []Channel{
		Channel{ID: "TestDM"},
	}
	ChannelList = []Channel{
		Channel{ID: "TestPublic"},
	}
	tests := []struct {
		name         string
		args         args
		wantChantype int
	}{
		{
			name:         "Test for expected DM GetChannel",
			args:         args{channel: "TestDM"},
			wantChantype: 1,
		},
		{
			name:         "Test for expected Public GetChannel",
			args:         args{channel: "TestPublic"},
			wantChantype: 2,
		},
		{
			name:         "Test for nonexistant GetChannel",
			args:         args{channel: "TestNone"},
			wantChantype: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotChantype := GetChannel(tt.args.channel); gotChantype != tt.wantChantype {
				t.Errorf("GetChannel() = %v, want %v", gotChantype, tt.wantChantype)
			}
		})
	}
}

func TestAddChannel(t *testing.T) {
	type args struct {
		channel  string
		chantype int
	}
	tests := []struct {
		name         string
		args         args
		wantChanList []Channel
	}{
		{
			name: "Test DM Channel",
			args: args{
				channel:  "test",
				chantype: 1,
			},
			wantChanList: DMChannelList,
		},
		{
			name: "Test Public Channel",
			args: args{
				channel:  "test",
				chantype: 2,
			},
			wantChanList: ChannelList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddChannel(tt.args.channel, tt.args.chantype)
			if checkArgIn(tt.args.channel, tt.wantChanList) {
				t.Errorf("AddChannel() = %v, want %v", false, tt.wantChanList)
			}
		})
	}
}

func checkArgIn(a string, list []Channel) bool {
	for _, b := range list {
		if b.ID == a {
			return true
		}
	}
	return false
}
