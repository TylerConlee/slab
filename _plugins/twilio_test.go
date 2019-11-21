package plugins

import (
	"reflect"
	"testing"

	"github.com/nlopes/slack"
)

func TestPlugins_EnableTwilio(t *testing.T) {
	type fields struct {
		Twilio    Twilio
		PagerDuty PagerDuty
	}
	tests := []struct {
		name        string
		fields      fields
		wantEnabled bool
	}{
		{
			name: "Check for Twilio Enabled",
			fields: fields{
				Twilio: Twilio{
					AccountID: "string",
					Auth:      "string",
				},
			},
			wantEnabled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugins{
				Twilio:    tt.fields.Twilio,
				PagerDuty: tt.fields.PagerDuty,
			}
			p.EnableTwilio()

			if TwilioEnabled != tt.wantEnabled {
				t.Errorf("Plugins.EnableTwilio() = %v, want %v", TwilioEnabled, tt.wantEnabled)
			}
		})
	}
}

func TestPlugins_DisableTwilio(t *testing.T) {
	type fields struct {
		Twilio    Twilio
		PagerDuty PagerDuty
	}
	tests := []struct {
		name         string
		fields       fields
		wantDisabled bool
	}{
		{
			name: "Check for Twilio Enabled",
			fields: fields{
				Twilio: Twilio{
					AccountID: "string",
					Auth:      "string",
				},
			},
			wantDisabled: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugins{
				Twilio:    tt.fields.Twilio,
				PagerDuty: tt.fields.PagerDuty,
			}
			p.DisableTwilio()

			if TwilioEnabled != tt.wantDisabled {
				t.Errorf("Plugins.DisableTwilio() = %v, want %v", TwilioEnabled, tt.wantDisabled)
			}
		})
	}
}

func TestTwilioSet(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name            string
		args            args
		wantTwilioPhone string
		wantAttachment  slack.Attachment
	}{
		{
			name: "Test for Setting Twilio Phone Number",
			args: args{
				n: "1234567890",
			},
			wantTwilioPhone: "1234567890",
			wantAttachment: slack.Attachment{
				Title: "Twilio 'To' Phone Number Set",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "Current Phone Number",
						Value: "1234567890",
						Short: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAttachment := TwilioSet(tt.args.n); !reflect.DeepEqual(gotAttachment, tt.wantAttachment) {
				t.Errorf("TwilioSet() = %v, want %v", gotAttachment, tt.wantAttachment)
			}
			if TwilioPhone != tt.wantTwilioPhone {
				t.Errorf("TwilioSet() = %v, want %v", TwilioPhone, tt.wantTwilioPhone)
			}
		})
	}
}

func TestTwilioUnset(t *testing.T) {
	tests := []struct {
		name            string
		wantTwilioPhone string
		wantAttachment  slack.Attachment
	}{
		{
			name:            "Test for Unsetting Twilio Phone Number",
			wantTwilioPhone: "",
			wantAttachment: slack.Attachment{
				Title: "Twilio 'To' Phone Number Unset",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "Current Phone Number",
						Value: "",
						Short: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAttachment := TwilioUnset(); !reflect.DeepEqual(gotAttachment, tt.wantAttachment) {
				t.Errorf("TwilioUnset() = %v, want %v", gotAttachment, tt.wantAttachment)
			}
			if TwilioPhone != tt.wantTwilioPhone {
				t.Errorf("TwilioSet() = %v, want %v", TwilioPhone, tt.wantTwilioPhone)
			}
		})
	}
}
