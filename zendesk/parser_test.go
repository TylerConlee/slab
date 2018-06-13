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
			name:       "Check for Level One",
			args:       []string{"platinum"},
			wantPriLvl: "LevelOne",
		},
		{
			name:       "Check for Level Two",
			args:       []string{"gold"},
			wantPriLvl: "LevelTwo",
		},
		{
			name:       "Check for Level Three",
			args:       []string{"silver"},
			wantPriLvl: "LevelThree",
		},
		{
			name:       "Check for Level Four",
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
