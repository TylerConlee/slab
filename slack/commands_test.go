package slack

import (
	"testing"

	"github.com/nlopes/slack"
)

func Test_parseDMCommand(t *testing.T) {
	type args struct {
		text string
		user string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseDMCommand(tt.args.text, tt.args.user)
		})
	}
}

func Test_parseCommand(t *testing.T) {
	type args struct {
		text    string
		user    *slack.User
		channel string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseCommand(tt.args.text, tt.args.user, tt.args.channel)
		})
	}
}
