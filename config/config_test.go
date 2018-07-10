package config

import (
	"reflect"
	"testing"
)

func Test_defaultConfig(t *testing.T) {
	tests := []struct {
		name       string
		wantConfig Config
	}{
		{
			name: "Test Default Config Generation",
			wantConfig: Config{
				Zendesk: Zendesk{
					APIKey: "",
					User:   "",
					URL:    "",
				},
				SLA: SLA{
					LevelOne: Level{
						Tag: "platinum",
					},
					LevelTwo: Level{
						Tag: "gold",
					},
					LevelThree: Level{
						Tag: "silver",
					},
					LevelFour: Level{
						Tag: "bronze",
					},
				},
				Metadata:      Metadata{},
				TriageEnabled: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConfig := defaultConfig(); !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("defaultConfig() = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
