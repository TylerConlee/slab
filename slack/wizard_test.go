package slack

import "testing"

func TestStartWizard(t *testing.T) {
	type args struct {
		user string
	}
	tests := []struct {
		name     string
		args     args
		wantUser string
		wantStep int
	}{
		{
			name: "Start Wizard with String",
			args: args{
				user: "test",
			},
			wantStep: 0,
			wantUser: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartWizard(tt.args.user)
			if activeUser.step != tt.wantStep {
				t.Errorf("activeUser.step = %v, want %v", activeUser.step, tt.wantStep)
			}
			if activeUser.user != tt.wantUser {
				t.Errorf("activeUser.user = %v, want %v", activeUser.user, tt.wantUser)
			}

		})
	}
}
