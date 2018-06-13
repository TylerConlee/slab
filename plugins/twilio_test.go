package plugins

import (
	"testing"
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
