package slack

import (
	"reflect"
	"testing"

	"github.com/tylerconlee/slack"
)

func TestWhoIsMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WhoIsMessage()
		})
	}
}

func TestSLAMessage(t *testing.T) {
	type args struct {
		n      string
		ticket Ticket
		color  string
		user   string
		uid    int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SLAMessage(tt.args.n, tt.args.ticket, tt.args.color, tt.args.user, tt.args.uid)
		})
	}
}

func TestDiagMessage(t *testing.T) {
	type args struct {
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
			DiagMessage(tt.args.user)
		})
	}
}

func TestNewTicketMessage(t *testing.T) {
	type args struct {
		tickets []Ticket
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewTicketMessage(tt.args.tickets)
		})
	}
}

func TestStatusMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StatusMessage()
		})
	}
}

func TestShowConfigMessage(t *testing.T) {
	type args struct {
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
			ShowConfigMessage(tt.args.user)
		})
	}
}

func TestHelpMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HelpMessage()
		})
	}
}

func TestUnknownCommandMessage(t *testing.T) {
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
			UnknownCommandMessage(tt.args.text, tt.args.user)
		})
	}
}

func TestSetMessage(t *testing.T) {
	tests := []struct {
		name           string
		wantAttachment slack.Attachment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAttachment := SetMessage(); !reflect.DeepEqual(gotAttachment, tt.wantAttachment) {
				t.Errorf("SetMessage() = %v, want %v", gotAttachment, tt.wantAttachment)
			}
		})
	}
}

func TestUnsetMessage(t *testing.T) {
	tests := []struct {
		name           string
		wantAttachment slack.Attachment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAttachment := UnsetMessage(); !reflect.DeepEqual(gotAttachment, tt.wantAttachment) {
				t.Errorf("UnsetMessage() = %v, want %v", gotAttachment, tt.wantAttachment)
			}
		})
	}
}
