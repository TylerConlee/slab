package zendesk

import "testing"

func Test_getPriorityLevel(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantPriLvl string
	}{
		// TODO: Add test cases.
		{
			name:       "LevelOne",
			args:       []string{"platinum"},
			wantPriLvl: "LevelOne",
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
