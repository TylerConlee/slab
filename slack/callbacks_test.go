package slack

import "testing"

func Test_statusDecode(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name    string
		args    args
		wantImg string
	}{
		{
			name: "Solved status decode",
			args: args{
				status: "solved",
			},
			wantImg: ":white_check_mark:",
		},
		{
			name: "New status decode",
			args: args{
				status: "new",
			},
			wantImg: ":new:",
		},
		{
			name: "Open status decode",
			args: args{
				status: "open",
			},
			wantImg: ":o2:",
		},
		{
			name: "Pending status decode",
			args: args{
				status: "pending",
			},
			wantImg: ":parking:",
		},
		{
			name: "Closed status decode",
			args: args{
				status: "closed",
			},
			wantImg: ":lock:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotImg := statusDecode(tt.args.status); gotImg != tt.wantImg {
				t.Errorf("statusDecode() = %v, want %v", gotImg, tt.wantImg)
			}
		})
	}
}
