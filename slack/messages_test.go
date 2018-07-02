package slack

import (
	"reflect"
	"testing"

	"github.com/tylerconlee/slack"
)

func TestSetMessage(t *testing.T) {
	Triager = "test"
	tests := []struct {
		name           string
		wantAttachment slack.Attachment
	}{
		{
			name: "Set Slab Triager",
			wantAttachment: slack.Attachment{
				Fallback:   "You would be able to select the triager here.",
				CallbackID: "triage_set",
				// Show the current triager
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "Current Triager",
						Value: "<@test>",
					},
				},

				// Show a dropdown of all users to select new Triager target
				Actions: []slack.AttachmentAction{
					slack.AttachmentAction{
						Name:  "triage_select",
						Text:  ":white_check_mark: Set",
						Type:  "button",
						Value: "ack",
						Style: "primary",
						Confirm: &slack.ConfirmationField{
							Text:        "Are you sure?",
							OkText:      "Take it",
							DismissText: "Leave it",
						},
					},
				},
			},
		},
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
		{
			name: "Unset Slab Triager",
			wantAttachment: slack.Attachment{
				Fallback:   "You would be able to select the triager here.",
				CallbackID: "triager_dropdown",
				Footer:     "Triager has been reset to None",
				FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAttachment := UnsetMessage(); !reflect.DeepEqual(gotAttachment, tt.wantAttachment) {
				t.Errorf("UnsetMessage() = %v, want %v", gotAttachment, tt.wantAttachment)
			}
		})
	}
}

func TestWhoIsMessage(t *testing.T) {
	Triager = "test"
	tests := []struct {
		name           string
		wantAttachment slack.Attachment
	}{
		{
			name: "Who Is Triager Message String",
			wantAttachment: slack.Attachment{
				Fallback:   "You would be able to select the triager here.",
				CallbackID: "triage_whois",
				// Show the current Triager member
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "Current Triager",
						Value: "<@test>",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAttachment := WhoIsMessage(); !reflect.DeepEqual(gotAttachment, tt.wantAttachment) {
				t.Errorf("WhoIsMessage() = %v, want %v", gotAttachment, tt.wantAttachment)
			}
		})
	}
}

func TestSLAMessage(t *testing.T) {
	type args struct {
		ticket Ticket
		color  string
		user   string
		uid    int64
	}
	tests := []struct {
		name           string
		args           args
		wantAttachment slack.Attachment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAttachment := SLAMessage(tt.args.ticket, tt.args.color, tt.args.user, tt.args.uid); !reflect.DeepEqual(gotAttachment, tt.wantAttachment) {
				t.Errorf("SLAMessage() = %v, want %v", gotAttachment, tt.wantAttachment)
			}
		})
	}
}
