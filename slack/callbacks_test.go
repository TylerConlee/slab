package slack

import (
	"testing"

	"github.com/tylerconlee/slab/datastore"
	"github.com/tylerconlee/slack"
)

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

func Test_satisfactionDecode(t *testing.T) {
	type args struct {
		sat string
	}
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		{
			name: "Good satisfaction decode",
			args: args{
				sat: "good",
			},
			wantS: ":white_check_mark:",
		},
		{
			name: "Bad satisfaction decode",
			args: args{
				sat: "bad",
			},
			wantS: ":x:",
		},
		{
			name: "Unoffered satisfaction decode",
			args: args{
				sat: "unoffered",
			},
			wantS: ":heavy_minus_sign:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := satisfactionDecode(tt.args.sat); gotS != tt.wantS {
				t.Errorf("satisfactionDecode() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestSetTriager(t *testing.T) {
	type args struct {
		payload *slack.AttachmentActionCallback
	}
	p := slack.AttachmentActionCallback{
		User: slack.User{
			ID: "test",
		},
	}
	tests := []struct {
		name        string
		args        args
		wantTriager string
	}{
		{
			name: "Save triager",
			args: args{
				payload: &p,
			},
			wantTriager: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTriager(tt.args.payload)
			if gotTriager := checkTriager(); gotTriager != tt.wantTriager {
				t.Errorf("checkTriager() = %v, want %v", gotTriager, tt.wantTriager)
			}
		})
	}
}

func checkTriager() (t string) {
	return datastore.Load("triager")
}
