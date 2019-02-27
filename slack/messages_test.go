package slack

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/nlopes/slack"
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

func TestSLAMessage(t *testing.T) {
	type Zendesk struct {
		URL string
	}
	type config struct {
		Zendesk Zendesk
	}

	type args struct {
		ticket Ticket
		color  string
		user   string
		uid    int64
		org    string
	}

	tests := []struct {
		name           string
		args           args
		wantAttachment slack.Attachment
	}{
		{
			name: "SLA Ticket Message",
			args: args{
				ticket: Ticket{
					Description: "This is a test ticket description",
					Priority:    "High",
					CreatedAt:   time.Now().Round(time.Second),
					Subject:     "Test Ticket",
					Requester:   123456,
					ID:          123,
				},
				color: "danger",
				user:  "test test",
				uid:   123456,
				org:   "lmnop",
			},
			wantAttachment: slack.Attachment{
				Title:      "Test Ticket",
				TitleLink:  "/agent/tickets/123",
				AuthorName: "test test",
				AuthorLink: "/agent/users/123456",
				AuthorIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/google/119/bust-in-silhouette_1f464.png",
				CallbackID: "sla",
				Color:      "danger",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "Description",
						Value: "This is a test ticket description",
					},
					slack.AttachmentField{
						Title: "Organization",
						Value: "lmnop",
					},
					slack.AttachmentField{
						Title: "Priority",
						Value: "High",
						Short: true,
					},
					slack.AttachmentField{
						Title: "Created At",
						Value: time.Now().Round(time.Second).String(),
						Short: true,
					},
				},
				Actions: []slack.AttachmentAction{
					slack.AttachmentAction{
						Name:  "ack_sla",
						Text:  ":white_check_mark: Acknowledge",
						Type:  "button",
						Value: "ack",
						Style: "primary",
						Confirm: &slack.ConfirmationField{
							Text:        "Are you sure?",
							OkText:      "Take it",
							DismissText: "Leave it",
						},
					},
					slack.AttachmentAction{
						Name:  "more_info_sla",
						Value: strconv.FormatInt(123456, 10),
						Text:  ":mag: More Info",
						Type:  "button",
						Style: "default",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAttachment := SLAMessage(tt.args.ticket, tt.args.color, tt.args.user, tt.args.uid, tt.args.org); !reflect.DeepEqual(gotAttachment, tt.wantAttachment) {
				t.Errorf("SLAMessage() = %v, want %v", gotAttachment, tt.wantAttachment)
			}
		})
	}
}
