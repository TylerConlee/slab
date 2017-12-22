package zendesk

import (
	"testing"
)

func Test_getPriorityLevel(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantPriLvl string
	}{
		{
			name:       "LevelOne",
			args:       []string{"platinum"},
			wantPriLvl: "LevelOne",
		},
		{
			name:       "LevelTwo",
			args:       []string{"gold"},
			wantPriLvl: "LevelTwo",
		},
		{
			name:       "LevelThree",
			args:       []string{"silver"},
			wantPriLvl: "LevelThree",
		},
		{
			name:       "LevelFour",
			args:       []string{"bronze"},
			wantPriLvl: "LevelFour",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPriLvl := getPriorityLevel(tt.args); gotPriLvl != tt.wantPriLvl {
				t.Errorf("getPriorityLevel() = %v, want %v", gotPriLvl, tt.wantPriLvl)
			}
		})
	}
}
